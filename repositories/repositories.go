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
	"sync"
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
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建 Hash 对象
	hasher := sha1.New()

	// 创建两个管道
	pipeReader, pipeWriter := io.Pipe()

	// 读取文件并写入到 pipeWriter 中
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pipeWriter.Close()
		_, err := io.Copy(pipeWriter, file)
		if err != nil {
			pipeWriter.CloseWithError(err)
		}
	}()

	// 读取 pipeReader 并计算 Hash
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(hasher, pipeReader)
		if err != nil {
			pipeReader.CloseWithError(err)
		}
	}()

	// 等待两个 Goroutine 完成
	wg.Wait()

	// 转换 Hash 成字符串
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
