package repositories

import (
	"NexusRepositorySync/config"
	"bufio"
	"errors"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type RepositoryFormat string

const (
	Maven2 RepositoryFormat = "maven2"
	Npm    RepositoryFormat = "npm"
)

const (
	HttpStatusCodeError = "ErrorCode"
	ConnectError        = "ConnectError"
)

type Repositoryer interface {
	Init()
	GetComponents(db *gorm.DB) error
	DownloadComponents(db *gorm.DB) error
	UploadComponents(db *gorm.DB) error
	Promote(s string)
}

type RepositoriesSync struct {
	DownloadRepository Repositoryer
	UploadRepository   Repositoryer
}

func httpGet(url, filePath, username, password string) error {
	method := "GET"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(username, password)
	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return errors.New(ConnectError)
	}
	defer res.Body.Close()
	if err := httpCodeCheck(res.StatusCode); err != nil {
		return err
	}
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	err = os.MkdirAll(filepath.Dir(filePath), 755)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	_, _ = io.Copy(writer, reader)
	return nil
}

func httpCodeCheck(statusCode int) error {
	if statusCode > 204 {
		log.Printf("返回状态码错误：%d\n", statusCode)
		return errors.New(HttpStatusCodeError)
	}
	return nil
}

func GetLocalFilePath(repositoryName string, assetPath string) string {
	localFilePath := filepath.Join(config.NexusConfig.DownloadPath, repositoryName, assetPath)
	return localFilePath
}
