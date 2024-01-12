package main

import (
	"NexusRepositorySync/config"
	"NexusRepositorySync/orm"
	"NexusRepositorySync/task"
	"NexusRepositorySync/web"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

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

func init() {
	err := os.MkdirAll(config.DownLoadDir, 755)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("任务执行间隔为：%v", task.TimeStep)
	log.Printf("监听端口为：%d", config.NexusConfig.Port)

	gin.SetMode(gin.ReleaseMode)
}

func main() {
	for _, repositorySync := range task.GetRepositorySyncTasks() {
		go task.RepositorySync(repositorySync, Db)
	}

	r := web.GetRouter()
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}

}
