package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const DownLoadDir = "download"

var NexusConfig Config

type RepositorySync struct {
	TaskName             string `yaml:"taskName"`
	DownRepositoryUrl    string `yaml:"downRepositoryUrl"`
	DownRepositoryName   string `yaml:"downRepositoryName"`
	DownRepositoryAuth   string `yaml:"downRepositoryAuth"`
	UploadRepositoryUrl  string `yaml:"uploadRepositoryUrl"`
	UploadRepositoryName string `yaml:"uploadRepositoryName"`
	UploadRepositoryAuth string `yaml:"uploadRepositoryAuth"`
	RepositoryType       string `yaml:"repositoryType"`
}

type Config struct {
	RepositorySyncTask []RepositorySync `yaml:"repositorySyncTask"`
	TimeStep           int              `yaml:"timeStep"`
	Port               int              `yaml:"port"`
	DbPath             string           `yaml:"dbPath"`
	DownloadPath       string           `yaml:"downloadPath"`
}

func init() {
	filePath := ".config.yaml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("无法读取配置文件 %s: %v", filePath, err)
	}

	err = yaml.Unmarshal([]byte(data), &NexusConfig)
	if err != nil {
		log.Fatalf("无法解析 YAML 数据: %v", err)
	}
}
