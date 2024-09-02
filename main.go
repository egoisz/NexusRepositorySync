package main

import (
	"NexusRepositorySync/task"
	"NexusRepositorySync/web"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"sync"
)

func init() {
	//err := os.MkdirAll(config.NexusConfig.DownloadPath, 755)
	//if err != nil {
	//	log.Panic(err)
	//}

	//log.Printf("任务执行间隔为：%v", task.TimeStep)
	//log.Printf("监听端口为：%d", config.NexusConfig.Port)

	gin.SetMode(gin.ReleaseMode)
}

func main() {
	// 子命令模式运行一次
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "download":
			fmt.Println("executing script for download")
			var wg sync.WaitGroup
			wg.Add(len(task.GetRepositorySyncTasks()))
			for _, repositorySync := range task.GetRepositorySyncTasks() {
				go task.RepositoryDownload(repositorySync, &wg)
			}
			wg.Wait()
		case "upload":
			fmt.Println("executing script for upload")
			var wg sync.WaitGroup
			wg.Add(len(task.GetRepositorySyncTasks()))
			for _, repositorySync := range task.GetRepositorySyncTasks() {
				go task.RepositoryUpload(repositorySync, &wg)
			}
			wg.Wait()
		default:
			fmt.Println("expected 'download' or 'upload' subcommands")
			os.Exit(1)
		}
		// 子命令执行完毕，退出程序
		os.Exit(0)
	}

	// 开始守护进程
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositorySync(repositorySync)
	}

	r := web.GetRouter()
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}

}
