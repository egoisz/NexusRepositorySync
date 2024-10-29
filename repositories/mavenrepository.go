package repositories

import (
	"NexusRepositorySync/orm"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type MavenRepository struct {
	Url  string
	Name string
	Type RepositoryFormat
	Auth string
}

func (r MavenRepository) GetComponents(db *gorm.DB) error {
	method := "GET"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	cTK := ""
	//var itermSlice []orm.Item
	for {
		url := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", r.Url, r.Name)
		if cTK != "" {
			url = fmt.Sprintf("%s&&continuationToken=%s", url, cTK)
		}
		//fmt.Println(url)
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			log.Printf("%s\n", err)
			continue
		}
		req.Header.Add("accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			// 超时打印日志继续，不退出
			r.Promote(err.Error())
			continue
			//return err
		}
		var t orm.NexusRequest
		err = json.NewDecoder(res.Body).Decode(&t)
		if err != nil {
			r.Promote(err.Error())
		}

		err = res.Body.Close()
		if err != nil {
			r.Promote(err.Error())
		}
		for _, item := range t.Items {
			for _, asset := range item.Assets {
				localFilePath := GetLocalFilePath(r.Name, asset.Path)
				db.Where(orm.MavenRepository{DownloadURL: asset.DownloadURL}).FirstOrCreate(&orm.MavenRepository{
					DownloadURL:   asset.DownloadURL,
					Path:          asset.Path,
					LocalFilePath: localFilePath,
					GroupID:       asset.Maven2.GroupID,
					ArtifactID:    asset.Maven2.ArtifactID,
					Version:       asset.Maven2.Version,
					Extension:     asset.Maven2.Extension,
					Classifier:    asset.Maven2.Classifier,
				})
				//err := HttpGet(v.DownloadURL, v.Path)
				if err != nil {
					r.Promote(err.Error())
				}
			}
		}

		cTK = t.ContinuationToken
		if cTK == "" {
			break
		}
	}
	return nil
}

func (r MavenRepository) DownloadComponents(db *gorm.DB) error {
	var t []orm.MavenRepository
	db.Where("down_load_status =?", false).Find(&t)
	//fmt.Println(n)
	for _, v := range t {
		//fmt.Println(v.DownloadURL, v.DownLoadStatus)
		//filePath := path.Join(config.DownLoadDir, r.Name, v.Path)
		err := httpGet(v.DownloadURL, v.LocalFilePath)
		if err != nil {
			r.Promote(fmt.Sprintf("下载失败：%s 原因：%s\n", v.DownloadURL, err.Error()))
			if err.Error() == HttpStatusCodeError {
				continue
			} else if err.Error() == ConnectError {
				return err
			}
			return err
		} else {
			r.Promote(fmt.Sprintf("下载完成 %s，%s\n", v.DownloadURL, v.LocalFilePath))
			db.Where(orm.MavenRepository{DownloadURL: v.DownloadURL}).Updates(orm.MavenRepository{DownLoadStatus: true})
		}
	}
	return nil
}

func (r MavenRepository) UploadComponents(db *gorm.DB) error {

	var n []orm.MavenRepository
	db.Where(
		"down_load_status =? and up_load_status =? and up_load_times <?", true, false, 3,
	).Where(
		"extension =? or extension=?", "pom", "jar").Find(&n)
	for _, v := range n {
		//url := http://10.147.235.204:8081/service/rest/v1/components?repository=inner-maven-public
		url := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", r.Url, r.Name)

		//auth := "Basic YWRtaW46WXl5dEBuZXh1c0AyMDIz"
		auth := fmt.Sprintf("Basic %s", r.Auth)
		//fmt.Println(v.DownloadURL, v.DownLoadStatus, v.Extension, v.UpLoadStatus)
		err := MavenComponentHttpPost(
			url,
			auth,
			v.LocalFilePath,
			v.GroupID,
			v.ArtifactID,
			v.Version,
			v.Extension,
			v.Classifier,
		)
		if err != nil {
			r.Promote(fmt.Sprintf("上传 %s 失败, 失败原因：%s\n", v.LocalFilePath, err))
			if err.Error() == HttpStatusCodeError {
				db.Where(orm.MavenRepository{DownloadURL: v.DownloadURL}).Updates(orm.MavenRepository{UpLoadTimes: v.UpLoadTimes + 1})
				continue
			} else if err.Error() == ConnectError {
				return err
			}
			return err
		}
		db.Where(orm.MavenRepository{DownloadURL: v.DownloadURL}).Updates(orm.MavenRepository{UpLoadStatus: true})
		r.Promote(fmt.Sprintf("上传成功成功：%s \n", v.LocalFilePath))
	}
	return nil
}

func (r MavenRepository) Promote(s string) {
	log.Printf("%-20s %-25s %s", r.Name, r.Url, s)
}

func MavenComponentHttpPost(
	url string,
	auth string,
	filePath string,
	groupId string,
	artifactId string,
	version string,
	extension string,
	classifier string,
) error {
	//url := "http://10.147.235.204:8081/service/rest/v1/components?repository=inner-maven-public"

	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("maven2.groupId", groupId)
	_ = writer.WriteField("maven2.artifactId", artifactId)
	_ = writer.WriteField("maven2.version", version)
	if classifier != "" {
		_ = writer.WriteField("maven2.asset1.classifier", classifier)
	}
	file, errFile4 := os.Open(filePath)
	defer file.Close()
	part4, errFile4 := writer.CreateFormFile("maven2.asset1", filepath.Base(filePath))
	_, errFile4 = io.Copy(part4, file)
	if errFile4 != nil {
		//fmt.Println(errFile4)
		return errFile4
	}
	_ = writer.WriteField("maven2.asset1.extension", extension)
	err := writer.Close()
	if err != nil {
		//fmt.Println(err)
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return err
	}
	req.Header.Add("accept", " application/json")
	req.Header.Add("Content-Type", "multipart/form-data")
	//req.Header.Add("Authorization", "Basic YWRtaW46WXl5dEBuZXh1c0AyMDIz")
	req.Header.Add("Authorization", auth)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New(ConnectError)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if string(body) != "" {
		//fmt.Println(err)
		log.Printf("return body: %s\n", string(body))

	}
	if err := httpCodeCheck(res.StatusCode); err != nil {
		return err
	}
	return nil
}
