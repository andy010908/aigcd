package main

import (
	//"context"
	"time"

	"aigcd/core/logger"
	"aigcd/mj"
	"aigcd/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const serverPort string = ":5000"

func main() {
	logger.InitLogger("/var/tmp/aigcd/logs/my.log", "debug")
	stop := make(chan struct{})
	//ctx := context.Background()

	go mj.DiscordBot()

	serv := gin.Default()
	//add cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT"}
	corsConfig.AllowHeaders = []string{"page", "count", "apikey", "Origin", "Content-Type", "Content-Length", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.MaxAge = 12 * time.Hour
	serv.Use(cors.New(corsConfig))
	routers.InitRouter(serv)
	serv.Run(serverPort)
	<-stop
}
