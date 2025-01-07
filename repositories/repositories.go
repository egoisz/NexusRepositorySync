package repositories

import (
	"NexusRepositorySync/config"
	"bufio"
	"crypto/sha1"
	"encoding/hex"
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

type Repositoryer interface {
	Init()
	GetComponents(db *gorm.DB, taskName string) error
	DownloadComponents(db *gorm.DB, taskName string) error
	UploadComponents(db *gorm.DB, taskName string) error
	Promote(s string)
}

type RepositoriesSync struct {
	DownloadRepository Repositoryer
	UploadRepository   Repositoryer
	TaskName           string
}

func httpGet(url, filePath, username, password string) error {
	method := "GET"
	client := http.Client{
		//Timeout: 5 * time.Second,
		Timeout: 0,
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

// CalculateFileSHA1 计算文件的 SHA1 值
func CalculateFileSHA1(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func checkMavenFileSuffix(s string) bool {
	for _, item := range config.NexusConfig.MavenFileSuffix {
		if item == s {
			return true
		}
	}
	return false
}
