package routers

import (
	"aigcd/controllers"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	r     *gin.Engine
	goapi *gin.RouterGroup
	rLock sync.Mutex
)

func InitRouter(r *gin.Engine) {
	rLock.Lock()
	defer rLock.Unlock()
	goapi = r.Group("")
	goapi.Use(gzip.Gzip(gzip.DefaultCompression))
	goapi.GET("/status", controllers.CheckStatus)
	//goapi.POST("/diffusion/send_prompt", controllers.SendPrompt)
	initMjRoute(goapi)
}
