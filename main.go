package main

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/orm"
	"NexusRepositorySync/repositories"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

const TimeStep = 10 * time.Second

var OutterMavenPublicRepository = repositories.MavenRepository{
	"http://172.30.86.46:18081",
	"maven-proxy-148-ali",
	repositories.Maven2,
	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
}

var UploadMavenPublicRepository = repositories.MavenRepository{
	"http://172.30.86.46:18081",
	"test-upload",
	repositories.Maven2,
	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
}

var OutterNpmPublicRepository = repositories.NpmRepository{
	"http://172.30.84.90:8081",
	"npm-local",
	repositories.Npm,
	"YWRtaW46WnlqY0AyMDIx",
}

var UploadNpmPublicRepository = repositories.NpmRepository{
	"http://172.30.86.46:18081",
	"test-npm-upload",
	repositories.Npm,
	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
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

func Syncrepository(repositorySyncsync repositories.RepositoriesSync, db *gorm.DB) {
	for {
		repositorySyncsync.DownloadRepository.Promote("开始获取组件清单")
		if err := repositorySyncsync.DownloadRepository.GetComponents(db); err != nil {
			repositorySyncsync.DownloadRepository.Promote(err.Error())
			repositorySyncsync.DownloadRepository.Promote("获取组件中断")
		} else {
			repositorySyncsync.DownloadRepository.Promote("获取组件清单结束")
		}

		repositorySyncsync.DownloadRepository.Promote("开始下载组件")
		if err := repositorySyncsync.DownloadRepository.DownloadComponents(db); err != nil {
			repositorySyncsync.DownloadRepository.Promote(err.Error())
			repositorySyncsync.DownloadRepository.Promote("下载组件中断")
		} else {
			repositorySyncsync.DownloadRepository.Promote("下载组件结束")
		}

		repositorySyncsync.UploadRepository.Promote("开始上传组件")
		if err := repositorySyncsync.UploadRepository.UploadComponents(db); err != nil {
			repositorySyncsync.UploadRepository.Promote(err.Error())
			repositorySyncsync.UploadRepository.Promote("上传组件中断")
		} else {
			repositorySyncsync.DownloadRepository.Promote("上传组件结束")
		}

		time.Sleep(TimeStep)
	}
}

func forever() {
	for {
		//An example goroutine that might run
		//indefinitely. In actual implementation
		//it might block on a chanel receive instead
		//of time.Sleep for example.
		//fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
func init() {
	err := os.MkdirAll(config.DownLoadDir, 755)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	var repositorySyncSice = []repositories.RepositoriesSync{
		{
			DownloadRepository: OutterMavenPublicRepository,
			UploadRepository:   UploadMavenPublicRepository,
		},
		{
			DownloadRepository: OutterNpmPublicRepository,
			UploadRepository:   UploadNpmPublicRepository,
		},
	}

	for _, repositorySyncsync := range repositorySyncSice {
		go Syncrepository(repositorySyncsync, Db)
	}

	forever()

	//if err := OutterMavenPublicRepository.GetComponents(Db); err != nil {
	//	fmt.Println(err)
	//}
	//if err := OutterMavenPublicRepository.DownloadComponents(Db); err != nil {
	//	fmt.Println(err)
	//}
	//if err := UploadMavenPublicRepository.UploadComponents(Db); err != nil {
	//	fmt.Println(err)
	//}
	//

	//if err := OutterNpmPublicRepository.GetComponents(Db); err != nil {
	//	fmt.Println(err)
	//}
	//if err := OutterNpmPublicRepository.DownloadComponents(Db); err != nil {
	//	fmt.Println(err)
	//}
	//if err := UploadNpmPublicRepository.UploadComponents(Db); err != nil {
	//	fmt.Println(err)
	//}
}
