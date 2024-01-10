package main

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/orm"
	"NexusRepositorySync/repositories"
	"NexusRepositorySync/web"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var TimeStep time.Duration
var Db = initDB()
var RepositorySyncTask []repositories.RepositoriesSync

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
		time.Sleep(time.Second)
	}
}
func init() {
	err := os.MkdirAll(config.DownLoadDir, 755)
	if err != nil {
		log.Panic(err)
	}
	for _, task := range config.NexusConfig.RepositorySyncTask {
		if task.RepositoryType == string(repositories.Maven2) {
			dRepository := repositories.MavenRepository{
				Url:  task.DownRepositoryUrl,
				Name: task.DownRepositoryName,
				Type: repositories.Maven2,
			}
			uRepository := repositories.MavenRepository{
				Url:  task.UploadRepositoryUrl,
				Name: task.UploadRepositoryName,
				Auth: task.UploadRepositoryAuth,
				Type: repositories.Maven2,
			}
			RepositorySyncTask = append(RepositorySyncTask, repositories.RepositoriesSync{
				DownloadRepository: dRepository,
				UploadRepository:   uRepository,
			})

		} else if task.RepositoryType == string(repositories.Npm) {
			dRepository := repositories.NpmRepository{
				Url:  task.DownRepositoryUrl,
				Name: task.DownRepositoryName,
				Type: repositories.Npm,
			}
			uRepository := repositories.NpmRepository{
				Url:  task.UploadRepositoryUrl,
				Name: task.UploadRepositoryName,
				Auth: task.UploadRepositoryAuth,
				Type: repositories.Npm,
			}
			RepositorySyncTask = append(RepositorySyncTask, repositories.RepositoriesSync{
				DownloadRepository: dRepository,
				UploadRepository:   uRepository,
			})
		}

	}
	log.Printf("任务执行间隔为：%v", TimeStep)
	log.Printf("监听端口为：%d", config.NexusConfig.Port)
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	for _, repositorySync := range RepositorySyncTask {
		go Syncrepository(repositorySync, Db)
	}

	r := web.GetRouter()
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}

	//forever()

}
