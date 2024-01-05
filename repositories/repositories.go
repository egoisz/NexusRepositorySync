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

// http://10.147.235.204:8081

type Repository interface {
	GetComponents(db *gorm.DB) error
	DownloadComponents(db *gorm.DB) error
	UploadComponents(db *gorm.DB) error
	Promote(s string)
}

type RepositoriesSync struct {
	DownloadRepository Repository
	UploadRepository   Repository
}

func httpGet(url string, filePath string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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
	localFilePath := filepath.Join(config.DownLoadDir, repositoryName, assetPath)
	return localFilePath
}
