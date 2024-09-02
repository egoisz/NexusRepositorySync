package main

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/task"
	"NexusRepositorySync/web"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func init() {
	//err := os.MkdirAll(config.NexusConfig.DownloadPath, 755)
	//if err != nil {
	//	log.Panic(err)
	//}

	log.Printf("任务执行间隔为：%v", task.TimeStep)
	log.Printf("监听端口为：%d", config.NexusConfig.Port)

	gin.SetMode(gin.ReleaseMode)
}

func main() {
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositorySync(repositorySync)
	}

	r := web.GetRouter()
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}

}
