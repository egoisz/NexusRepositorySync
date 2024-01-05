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

//var OutterMavenPublicRepository = repositories.MavenRepository{
//	"http://172.30.86.46:18081",
//	"maven-proxy-148-ali",
//	repositories.Maven2,
//	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
//}
//
//var UploadMavenPublicRepository = repositories.MavenRepository{
//	"http://172.30.86.46:18081",
//	"test-upload",
//	repositories.Maven2,
//	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
//}
//
//var OutterNpmPublicRepository = repositories.NpmRepository{
//	"http://172.30.84.90:8081",
//	"npm-local",
//	repositories.Npm,
//	"YWRtaW46WnlqY0AyMDIx",
//}
//
//var UploadNpmPublicRepository = repositories.NpmRepository{
//	"http://172.30.86.46:18081",
//	"test-npm-upload",
//	repositories.Npm,
//	"YWRtaW46SHlkZXZAbmV4dXMyMDIz",
//}

// TransMavenPublicRepository prod 配置
var TransMavenPublicRepository = repositories.MavenRepository{
	Url:  "http://172.30.86.46:18081",
	Name: "sync-maven-public",
	Type: repositories.Maven2,
	Auth: "YWRtaW46SHlkZXZAbmV4dXMyMDIz",
}

var TransNpmPublicRepository = repositories.NpmRepository{
	Url:  "http://172.30.86.46:18081",
	Name: "sync-npm-public",
	Type: repositories.Npm,
	Auth: "YWRtaW46SHlkZXZAbmV4dXMyMDIz",
}

var InnerMavenPublicRepository = repositories.MavenRepository{
	Url:  "http://10.147.235.204:8081",
	Name: "inner-maven-public",
	Type: repositories.Maven2,
	Auth: "YWRtaW46WXl5dEBuZXh1c0AyMDIz",
}

var InnerNpmPublicRepository = repositories.NpmRepository{
	Url:  "http://10.147.235.204:8081",
	Name: "inner-npm-public",
	Type: repositories.Npm,
	Auth: "YWRtaW46WXl5dEBuZXh1c0AyMDIz",
}

var Db = initDB()

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("nexus.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&orm.MavenRepository{})
	_ = db.AutoMigrate(&orm.NpmRepository{})
	return db
}

func Syncrepository(r repositories.RepositoriesSync, db *gorm.DB) {
	for {
		r.DownloadRepository.Promote("开始获取组件清单")
		if err := r.DownloadRepository.GetComponents(db); err != nil {
			r.DownloadRepository.Promote(err.Error())
			r.DownloadRepository.Promote("获取组件中断")
		} else {
			r.DownloadRepository.Promote("获取组件清单结束")
		}

		r.DownloadRepository.Promote("开始下载组件")
		if err := r.DownloadRepository.DownloadComponents(db); err != nil {
			r.DownloadRepository.Promote(err.Error())
			r.DownloadRepository.Promote("下载组件中断")
		} else {
			r.DownloadRepository.Promote("下载组件结束")
		}

		r.UploadRepository.Promote("开始上传组件")
		if err := r.UploadRepository.UploadComponents(db); err != nil {
			r.UploadRepository.Promote(err.Error())
			r.UploadRepository.Promote("上传组件中断")
		} else {
			r.DownloadRepository.Promote("上传组件结束")
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
	// prod 使用
	var repositorySyncSice = []repositories.RepositoriesSync{
		{
			DownloadRepository: TransMavenPublicRepository,
			UploadRepository:   InnerMavenPublicRepository,
		},
		{
			DownloadRepository: TransNpmPublicRepository,
			UploadRepository:   InnerNpmPublicRepository,
		},
	}

	for _, repositorySync := range repositorySyncSice {
		go Syncrepository(repositorySync, Db)
	}

	forever()

	// 测试
	//var repositorySyncSice = []repositories.RepositoriesSync{
	//	{
	//		DownloadRepository: OutterMavenPublicRepository,
	//		UploadRepository:   InnerMavenPublicRepository,
	//	},
	//	{
	//		DownloadRepository: OutterNpmPublicRepository,
	//		UploadRepository:   InnerNpmPublicRepository,
	//	},
	//}

	// dev
	//var repositorySyncSice = []repositories.RepositoriesSync{
	//	{
	//		DownloadRepository: OutterMavenPublicRepository,
	//		UploadRepository:   UploadMavenPublicRepository,
	//	},
	//	{
	//		DownloadRepository: OutterNpmPublicRepository,
	//		UploadRepository:   UploadNpmPublicRepository,
	//	},
	//}
	//
	//for _, repositorySyncsync := range repositorySyncSice {
	//	go Syncrepository(repositorySyncsync, Db)
	//}

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
