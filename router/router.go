package router

import (
	"BlueNetDisk/api"
	"BlueNetDisk/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	v1 := ginRouter.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.POST("user/register", api.UserRegisterHandler())
		v1.POST("user/login", api.UserLoginHandler())

		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("authed/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "authed pong!",
				})
			})
			authed.POST("file/upload", api.UploadFileHandler())
			authed.GET("file/list", api.ListFileHandler())
			authed.POST("file/mkdir", api.MkdirHandler())
		}
	}
	return ginRouter
}
