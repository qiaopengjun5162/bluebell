package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)

	r.GET("/version", func(context *gin.Context) {
		context.String(http.StatusOK, setting.Conf.Version)
	})
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "404",
		})
	})
	return r
}
