package main

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/orm"
	"NexusRepositorySync/repositories"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

var OutterMavenPublicRepository = repositories.MavenRepository{
	"http://172.30.86.46:18081",
	"maven-proxy-148-ali",
	repositories.Maven2,
	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
}

var OutterNpmPublicRepository = repositories.NpmRepository{
	"http://172.30.84.90:8081",
	"npm-local",
	repositories.Npm,
	"YWRtaW46WnlqY0AyMDIx",
}

var Db = initDB()

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("nexus.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&orm.MavenRepository{})
	db.AutoMigrate(&orm.NpmRepository{})
	return db
}

func init() {
	err := os.MkdirAll(config.DownLoadDir, 755)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	//err := OutterMavenPublicRepository.GetComponents(Db)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = OutterMavenPublicRepository.DownloadComponents(Db)
	//if err != nil {
	//	fmt.Println(err)
	//}
	err := OutterNpmPublicRepository.GetComponents(Db)
	if err != nil {
		fmt.Println(err)
	}
	err = OutterNpmPublicRepository.DownloadComponents(Db)
	if err != nil {
		fmt.Println(err)
	}
}
