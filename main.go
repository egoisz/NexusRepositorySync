package main

import (
	"NexusRepositorySync/task"
	"NexusRepositorySync/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sync"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func executeSearch() {
	fmt.Println("executing script for search")
	var wg sync.WaitGroup
	wg.Add(len(task.GetRepositorySyncTasks()))
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositorySearch(repositorySync, &wg)
	}
	wg.Wait()
}

func executeDownload() {
	fmt.Println("executing script for download")
	var wg sync.WaitGroup
	wg.Add(len(task.GetRepositorySyncTasks()))
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositoryDownload(repositorySync, &wg)
	}
	wg.Wait()
}

func executeUpload() {
	fmt.Println("executing script for upload")
	var wg sync.WaitGroup
	wg.Add(len(task.GetRepositorySyncTasks()))
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositoryUpload(repositorySync, &wg)
	}
	wg.Wait()
}

func startDaemon() {
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositorySync(repositorySync)
	}

	r := web.GetRouter()
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// 子命令模式运行一次
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "search":
			executeSearch()
		case "download":
			executeDownload()
		case "upload":
			executeUpload()
		default:
			fmt.Println("expected 'search' or 'download' or 'upload' subcommands")
			os.Exit(1)
		}
		// 子命令执行完毕，退出程序
		os.Exit(0)
	}

	// 开始守护进程
	startDaemon()
}
