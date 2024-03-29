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

type NpmRepository struct {
	Url  string
	Name string
	Type RepositoryFormat
	Auth string
}

func (r NpmRepository) GetComponents(db *gorm.DB) error {
	method := "GET"
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	cTK := ""

	for {
		url := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", r.Url, r.Name)
		if cTK != "" {
			url = fmt.Sprintf("%s&&continuationToken=%s", url, cTK)
		}
		//fmt.Println(url)
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			r.Promote(err.Error())
			continue
		}
		req.Header.Add("accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			//r.Promote(err.Error())
			return err
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
				db.Where(orm.NpmRepository{DownloadURL: asset.DownloadURL}).FirstOrCreate(&orm.NpmRepository{
					DownloadURL:   asset.DownloadURL,
					Name:          asset.Npm.Name,
					LocalFilePath: localFilePath,
					Path:          asset.Path,
					Version:       asset.Npm.Version,
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

func (r NpmRepository) DownloadComponents(db *gorm.DB) error {
	var t []orm.NpmRepository
	db.Where("down_load_status =?", false).Find(&t)
	//fmt.Println(n)
	for _, v := range t {
		//fmt.Println(v.DownloadURL, v.DownLoadStatus)
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
			db.Where(orm.NpmRepository{DownloadURL: v.DownloadURL}).Updates(orm.NpmRepository{DownLoadStatus: true})
		}
	}
	return nil
}

func (r NpmRepository) UploadComponents(db *gorm.DB) error {

	var n []orm.NpmRepository
	db.Where(
		"down_load_status =? and up_load_status =?", true, false,
	).Find(&n)
	for _, v := range n {
		//url := http://10.147.235.204:8081/service/rest/v1/components?repository=inner-maven-public
		url := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", r.Url, r.Name)

		//auth := "Basic YWRtaW46WXl5dEBuZXh1c0AyMDIz"
		auth := fmt.Sprintf("Basic %s", r.Auth)
		//fmt.Println(v.DownloadURL, v.DownLoadStatus, v.Extension, v.UpLoadStatus)
		err := NpmComponentHttpPost(
			url,
			auth,
			v.LocalFilePath)
		if err != nil {
			r.Promote(fmt.Sprintf("上传 %s 失败, 失败原因：%s\n,", v.LocalFilePath, err))
			if err.Error() == HttpStatusCodeError {
				continue
			} else if err.Error() == ConnectError {
				return err
			}
			return err
		}
		db.Where(orm.NpmRepository{DownloadURL: v.DownloadURL}).Updates(orm.NpmRepository{UpLoadStatus: true})
		r.Promote(fmt.Sprintf("上传成功成功：%s \n", v.LocalFilePath))
	}
	return nil
}

func (r NpmRepository) Promote(s string) {
	log.Printf("%-20s %-25s %s", r.Name, r.Url, s)
}

func NpmComponentHttpPost(
	url string,
	auth string,
	filePath string,
) error {

	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile4 := os.Open(filePath)
	defer file.Close()
	part4, errFile4 := writer.CreateFormFile("npm.asset1", filepath.Base(filePath))
	_, errFile4 = io.Copy(part4, file)
	if errFile4 != nil {
		//fmt.Println(errFile4)
		return errFile4
	}
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
