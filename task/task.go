package task

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/repositories"
	"gorm.io/gorm"
	"time"
)

var TimeStep = time.Duration(config.NexusConfig.TimeStep) * time.Second

func RepositorySync(r repositories.RepositoriesSync, db *gorm.DB) {
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

func GetRepositorySyncTasks() []repositories.RepositoriesSync {
	var repositorySyncTask []repositories.RepositoriesSync
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
			repositorySyncTask = append(repositorySyncTask, repositories.RepositoriesSync{
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
			repositorySyncTask = append(repositorySyncTask, repositories.RepositoriesSync{
				DownloadRepository: dRepository,
				UploadRepository:   uRepository,
			})
		}

	}
	return repositorySyncTask
}
