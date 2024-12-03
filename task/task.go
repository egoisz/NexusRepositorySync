package task

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/orm"
	"NexusRepositorySync/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"sync"
	"time"
)

var TimeStep = time.Duration(config.NexusConfig.TimeStep) * time.Second

func initDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&orm.MavenRepository{})
	_ = db.AutoMigrate(&orm.NpmRepository{})
	return db
}

func RepositorySearch(r repositories.RepositoriesSync, wg *sync.WaitGroup) {
	defer wg.Done()
	dbPath := config.NexusConfig.DbPath
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		r.DownloadRepository.Promote("数据库文件不存在,跳过本次同步")
		return
	}
	db := initDB(dbPath)

	r.DownloadRepository.Promote("开始获取组件清单")
	if err := r.DownloadRepository.GetComponents(db); err != nil {
		r.DownloadRepository.Promote(err.Error())
		r.DownloadRepository.Promote("获取组件中断")
	} else {
		r.DownloadRepository.Promote("获取组件清单结束")
	}
}

func RepositoryDownload(r repositories.RepositoriesSync, wg *sync.WaitGroup) {
	defer wg.Done()
	dbPath := config.NexusConfig.DbPath
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		r.DownloadRepository.Promote("数据库文件不存在,跳过本次同步")
		return
	}
	db := initDB(dbPath)

	r.DownloadRepository.Promote("开始下载组件")
	if err := r.DownloadRepository.DownloadComponents(db); err != nil {
		r.DownloadRepository.Promote(err.Error())
		r.DownloadRepository.Promote("下载组件中断")
	} else {
		r.DownloadRepository.Promote("下载组件结束")
	}
}

func RepositoryUpload(r repositories.RepositoriesSync, wg *sync.WaitGroup) {
	defer wg.Done()
	dbPath := config.NexusConfig.DbPath
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		r.DownloadRepository.Promote("数据库文件不存在,跳过本次同步")
		return
	}
	db := initDB(dbPath)
	r.UploadRepository.Promote("开始上传组件")
	if err := r.UploadRepository.UploadComponents(db); err != nil {
		r.UploadRepository.Promote(err.Error())
		r.UploadRepository.Promote("上传组件中断")
	} else {
		r.DownloadRepository.Promote("上传组件结束")
	}

}

func RepositorySync(r repositories.RepositoriesSync) {
	for {
		dbPath := config.NexusConfig.DbPath
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			r.DownloadRepository.Promote("数据库文件不存在,跳过本次同步")
			time.Sleep(TimeStep)
			continue
		}
		db := initDB(dbPath)

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

func GetRepositorySyncTasks() []repositories.RepositoriesSync {
	var repositorySyncTask []repositories.RepositoriesSync
	for _, task := range config.NexusConfig.RepositorySyncTask {
		var dRepository, uRepository repositories.Repositoryer
		switch task.RepositoryType {
		case string(repositories.Maven2):
			dRepository = &repositories.MavenRepository{
				Url:  task.DownRepositoryUrl,
				Name: task.DownRepositoryName,
				Auth: task.DownRepositoryAuth,
				Type: repositories.Maven2,
			}
			uRepository = &repositories.MavenRepository{
				Url:  task.UploadRepositoryUrl,
				Name: task.UploadRepositoryName,
				Auth: task.UploadRepositoryAuth,
				Type: repositories.Maven2,
			}

		case string(repositories.Npm):
			dRepository = &repositories.NpmRepository{
				Url:  task.DownRepositoryUrl,
				Name: task.DownRepositoryName,
				Auth: task.DownRepositoryAuth,
				Type: repositories.Npm,
			}
			uRepository = &repositories.NpmRepository{
				Url:  task.UploadRepositoryUrl,
				Name: task.UploadRepositoryName,
				Auth: task.UploadRepositoryAuth,
				Type: repositories.Npm,
			}
		default:
			continue // 如果是未知类型，跳过
		}

		dRepository.Init()
		uRepository.Init()
		repositorySyncTask = append(repositorySyncTask, repositories.RepositoriesSync{
			DownloadRepository: dRepository,
			UploadRepository:   uRepository,
		})
	}

	return repositorySyncTask
}
