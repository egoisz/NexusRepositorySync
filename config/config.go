package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const DownLoadDir = "download"

var NexusConfig Config

type Config struct {
	DownloadMavenRepositoryUrl  string `yaml:"downloadMavenRepositoryUrl"`
	DownloadMavenRepositoryName string `yaml:"downloadMavenRepositoryName"`
	DownloadMavenRepositoryAuth string `yaml:"downloadMavenRepositoryAuth"`
	UploadMavenRepositoryUrl    string `yaml:"uploadMavenRepositoryUrl"`
	UploadMavenRepositoryName   string `yaml:"uploadMavenRepositoryName"`
	UploadMavenRepositoryAuth   string `yaml:"uploadMavenRepositoryAuth"`
	DownloadNpmRepositoryUrl    string `yaml:"downloadNpmRepositoryUrl"`
	DownloadNpmRepositoryName   string `yaml:"downloadNpmRepositoryName"`
	DownloadNpmRepositoryAuth   string `yaml:"downloadNpmRepositoryAuth"`
	UploadNpmRepositoryUrl      string `yaml:"uploadNpmRepositoryUrl"`
	UploadNpmRepositoryName     string `yaml:"uploadNpmRepositoryName"`
	UploadNpmRepositoryAuth     string `yaml:"uploadNpmRepositoryAuth"`
	TimeStep                    int    `yaml:"timeStep"`
	Port                        int    `yaml:"port"`
}

func init() {
	filePath := ".config.yaml" // 将此处修改为你自己的 YAML 文件路径
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("无法读取配置文件 %s: %v", filePath, err)
	}

	err = yaml.Unmarshal([]byte(data), &NexusConfig)
	if err != nil {
		log.Fatalf("无法解析 YAML 数据: %v", err)
	}

	//fmt.Println("从 YAML 文件中读取到的内容:")
	//fmt.Printf("姓名: %s\n年龄: %d\n", config.Name, config.Age)
}
