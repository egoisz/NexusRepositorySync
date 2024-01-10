package web

import (
	"NexusRepositorySync/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRouter() *gin.Engine {
	var r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	http.ListenAndServe(fmt.Sprintf(":%d", config.NexusConfig.Port), r)
	return r
}
