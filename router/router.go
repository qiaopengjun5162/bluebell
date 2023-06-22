package router

import (
	"bluebell/logger"
	"bluebell/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/version", func(context *gin.Context) {
		context.String(http.StatusOK, setting.Conf.Version)
	})
	return r
}
