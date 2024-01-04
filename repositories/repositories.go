package repositories

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/orm"
	"bufio"
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
	"path"
	"path/filepath"
	"time"
)

type RepositoryFormat string

const (
	Maven2 RepositoryFormat = "maven2"
	Npm    RepositoryFormat = "npm"
)

// http://10.147.235.204:8081

type Repository interface {
	GetComponents(db *gorm.DB) error
	DownloadComponents(db *gorm.DB) error
}

type MavenRepository struct {
	Url  string
	Name string
	Type RepositoryFormat
	Auth string
}

type RepositoriesSync struct {
	OuterRepository Repository
	InnerRepository Repository
}

func (r *MavenRepository) GetComponents(db *gorm.DB) error {
	method := "GET"
	client := &http.Client{
		Timeout: 5 * time.Second,
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
			log.Printf("%s\n", err)
			return err
		}
		var t orm.NexusRequest
		err = json.NewDecoder(res.Body).Decode(&t)
		if err != nil {
			log.Printf("%s\n", err)
		}

		err = res.Body.Close()
		if err != nil {
			log.Println(err)
		}
		for _, item := range t.Items {
			//fmt.Println(v.DownloadURL)
			if r.Type == Maven2 {
				for _, asset := range item.Assets {
					db.Where(orm.MavenRepository{DownloadURL: asset.DownloadURL}).FirstOrCreate(&orm.MavenRepository{
						DownloadURL: asset.DownloadURL,
						Path:        asset.Path,
						GroupID:     asset.Maven2.GroupID,
						ArtifactID:  asset.Maven2.ArtifactID,
						Version:     asset.Maven2.Version,
						Extension:   asset.Maven2.Extension,
					})
					//err := HttpGet(v.DownloadURL, v.Path)
					if err != nil {
						log.Printf("%s\n", err)
					}
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

func (r *MavenRepository) DownloadComponents(db *gorm.DB) error {
	var t []orm.MavenRepository
	db.Where("down_load_status =?", false).Find(&t)
	//fmt.Println(n)
	for _, v := range t {
		//fmt.Println(v.DownloadURL, v.DownLoadStatus)
		filePath := path.Join(config.DownLoadDir, r.Name, v.Path)
		err := httpGet(v.DownloadURL, filePath)
		if err != nil {
			log.Printf("下载失败%s\n", v.DownloadURL)
		}
		db.Where(orm.MavenRepository{DownloadURL: v.DownloadURL}).Updates(orm.MavenRepository{DownLoadStatus: true})
	}
	log.Printf("下载文件至本地完成!\n")
	return nil
}

func (r *MavenRepository) UploadComponents(db *gorm.DB) error {
	var t []orm.MavenRepository
	db.Where("down_load_status =?", false).Find(&t)
	//fmt.Println(n)
	for _, v := range t {
		//fmt.Println(v.DownloadURL, v.DownLoadStatus)
		filePath := path.Join(config.DownLoadDir, r.Name, v.Path)
		err := httpGet(v.DownloadURL, filePath)
		if err != nil {
			log.Printf("下载失败%s\n", v.DownloadURL)
		}
		db.Where(orm.MavenRepository{DownloadURL: v.DownloadURL}).Updates(orm.MavenRepository{DownLoadStatus: true})
	}
	log.Printf("下载文件至本地完成!\n")
	return nil
}

type NpmRepository struct {
	Url  string
	Name string
	Type RepositoryFormat
	Auth string
}

func (r *NpmRepository) GetComponents(db *gorm.DB) error {
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
			log.Printf("%s\n", err)
			continue
		}
		req.Header.Add("accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			log.Printf("%s\n", err)
			return err
		}
		var t orm.NexusRequest
		err = json.NewDecoder(res.Body).Decode(&t)
		if err != nil {
			log.Printf("%s\n", err)
		}

		err = res.Body.Close()
		if err != nil {
			log.Println(err)
		}
		for _, item := range t.Items {
			//fmt.Println(item)
			for _, asset := range item.Assets {
				db.Where(orm.NpmRepository{DownloadURL: asset.DownloadURL}).FirstOrCreate(&orm.NpmRepository{
					DownloadURL: asset.DownloadURL,
					Name:        asset.Npm.Name,
					Path:        asset.Path,
					Version:     asset.Npm.Version,
				})
				//err := HttpGet(v.DownloadURL, v.Path)
				if err != nil {
					log.Printf("%s\n", err)
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

func (r *NpmRepository) DownloadComponents(db *gorm.DB) error {
	var t []orm.NpmRepository
	db.Where("down_load_status =?", false).Find(&t)
	//fmt.Println(n)
	for _, v := range t {
		//fmt.Println(v.DownloadURL, v.DownLoadStatus)
		filePath := path.Join(config.DownLoadDir, r.Name, v.Path)
		err := httpGet(v.DownloadURL, filePath)
		if err != nil {
			log.Printf("下载失败%s\n", v.DownloadURL)
		}
		db.Where(orm.NpmRepository{DownloadURL: v.DownloadURL}).Updates(orm.NpmRepository{DownLoadStatus: true})
	}
	//log.Printf("下载文件至本地完成!\n")
	return nil
}

func httpGet(url string, filePath string) error {

	fileName := path.Base(url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	err = os.MkdirAll(filePath, 755)
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(filePath, fileName))
	if err != nil {
		return err
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	log.Printf("下载完成%s %s, Total length: %d \n", url, filePath, written)
	return nil
}

func httpPost(groupId string,
	artifactId string,
	version string,
	filePath string,
	extension string,
) error {
	//url := "http://172.30.86.46:18081/service/rest/v1/components?repository=test-upload"
	url := "http://10.147.235.204:8081/service/rest/v1/components?repository=inner-maven-public"

	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("maven2.groupId", groupId)
	_ = writer.WriteField("maven2.artifactId", artifactId)
	_ = writer.WriteField("maven2.version", version)
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
	//req.Header.Add("Authorization", "Basic YWRtaW46SHlkZXZAbmV4dXMyMDIz")
	req.Header.Add("Authorization", "Basic YWRtaW46WXl5dEBuZXh1c0AyMDIz")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("ConnetError")
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
	log.Println(res.StatusCode)
	return nil
}
