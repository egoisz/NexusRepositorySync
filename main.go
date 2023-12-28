package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
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

type NexusRequest struct {
	Items             []Item `json:"items"`
	ContinuationToken string `json:"continuationToken"`
}

type Item struct {
	ID         string  `json:"id"`
	Repository string  `json:"repository"`
	Format     string  `json:"format"`
	Group      string  `json:"group"`
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	Assets     []Asset `json:"assets"`
}

type Asset struct {
	DownloadURL    string      `json:"downloadUrl"`
	Path           string      `json:"path"`
	ID             string      `json:"id"`
	Repository     string      `json:"repository"`
	Format         string      `json:"format"`
	Checksum       Checksum    `json:"checksum"`
	ContentType    string      `json:"contentType"`
	LastModified   string      `json:"lastModified"`
	LastDownloaded string      `json:"lastDownloaded"`
	Uploader       string      `json:"uploader"`
	UploaderIP     string      `json:"uploaderIp"`
	FileSize       int64       `json:"fileSize"`
	BlobCreated    string      `json:"blobCreated"`
	Maven2         Maven2Class `json:"maven2"`
}

type Checksum struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
	Sha256 string `json:"sha256"`
	Md5    string `json:"md5"`
}

type Maven2Class struct {
	Extension  string  `json:"extension"`
	GroupID    string  `json:"groupId"`
	ArtifactID string  `json:"artifactId"`
	Version    string  `json:"version"`
	Classifier *string `json:"classifier,omitempty"`
}

type Nexus struct {
	gorm.Model
	DownloadURL    string
	GroupID        string
	ArtifactID     string
	Version        string
	Path           string
	Extension      string
	DownLoadStatus bool `gorm:"default:false"`
	UpLoadStatus   bool `gorm:"default:false"`
}

const DownLoadDir = "download"

var DB = initDB()

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("nexus.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func UploadNexusData() {
	for {
		time.Sleep(20 * time.Second)
		log.Printf("开始上传文件至nexus\n")
		var n []Nexus
		DB.Where(
			"down_load_status =? and up_load_status =?", true, false,
		).Where(
			"extension =? or extension=?", "pom", "jar").Find(&n)
		for _, v := range n {
			//fmt.Println(v.DownloadURL, v.DownLoadStatus, v.Extension, v.UpLoadStatus)
			err := HttpPost(
				v.GroupID,
				v.ArtifactID,
				v.Version,
				path.Join(DownLoadDir, v.Path),
				v.Extension)
			if err != nil {
				log.Printf("上传 %s 失败: %s\n", v.DownloadURL, err)
				continue
			}
			DB.Where(Nexus{DownloadURL: v.DownloadURL}).Updates(Nexus{UpLoadStatus: true})
			log.Printf("上传成功成功%s \n", v.DownloadURL)
		}
		log.Printf("上传文件至nexus完成!\n")
	}
}

func UpdateNexusData() {

	method := "GET"
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	cTK := ""
	// 无线循环
	for {
		// 隔10秒请求一次
		time.Sleep(10 * time.Second)
		log.Printf("开始更新数据库\n")
		for {
			url := "http://172.30.86.46:18081/service/rest/v1/components?repository=maven-public"
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
				continue
			}
			var t NexusRequest
			err = json.NewDecoder(res.Body).Decode(&t)
			if err != nil {
				log.Printf("%s\n", err)
			}

			//db.Create(&Nexus{DownloadURL: v.DownloadURL, GroupID: v.Maven2.GroupID, ArtifactID: v.Maven2.ArtifactID, Version: v.Maven2.Version})
			//db.Create(&Nexus{DownloadURL: "v.DownloadURL"})
			err = res.Body.Close()
			if err != nil {
				log.Println(err)
			}
			for _, item := range t.Items {
				//fmt.Println(v.DownloadURL)
				// 不存在URL则新建目录
				for _, asset := range item.Assets {
					DB.Where(Nexus{DownloadURL: asset.DownloadURL}).FirstOrCreate(&Nexus{
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
			cTK = t.ContinuationToken
			if cTK == "" {
				break
			}
		}
		log.Printf("数据库更新完成!\n")
	}
}

func HttpPost(groupId string,
	artifactId string,
	version string,
	filePath string,
	extension string,
) error {
	url := "http://172.30.86.46:18081/service/rest/v1/components?repository=test-upload"
	//url := "http://10.147.235.204:8081/service/rest/v1/components?repository=maven-public"

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
	req.Header.Add("Authorization", "Basic YWRtaW46SHlkZXZAbmV4dXMyMDIz")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	if string(body) != "" {
		//fmt.Println(err)
		log.Printf("return body: %s\n", string(body))

	}
	return nil
}

func HttpGet(url string, pa string) error {

	fileName := path.Base(url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	filePath := path.Join(DownLoadDir, path.Dir(pa))

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
	log.Printf("下载完成%s %s, Total length: %d \n", url, pa, written)
	return nil
}
func init() {
	err := os.MkdirAll(DownLoadDir, 755)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	DB.AutoMigrate(&Nexus{})

	// 上传nexus仓库
	go UploadNexusData()
	// 更新数据库
	go UpdateNexusData()

	// 下载http文件
	for {
		time.Sleep(10 * time.Second)
		log.Printf("开始下载文件至本地\n")
		var n []Nexus
		DB.Where("down_load_status =?", false).Find(&n)
		//fmt.Println(n)
		for _, v := range n {
			fmt.Println(v.DownloadURL, v.DownLoadStatus)
			err := HttpGet(v.DownloadURL, v.Path)
			if err != nil {
				fmt.Printf("下载失败%s\n", v.DownloadURL)
			}
			DB.Where(Nexus{DownloadURL: v.DownloadURL}).Updates(Nexus{DownLoadStatus: true})
		}
		log.Printf("下载文件至本地完成!\n")
	}
}
