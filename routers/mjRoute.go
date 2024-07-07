package routers

import (
	"aigcd/controllers"

	"github.com/gin-gonic/gin"
)

func initMjRoute(v1 *gin.RouterGroup) {
	v1.POST("/diffusion/send_prompt", controllers.SendPrompt)
	v1.GET("/diffusion/collections", controllers.GetCollections)
	v1.POST("/diffusion/cloud_upload", controllers.CloudUpload)
}
