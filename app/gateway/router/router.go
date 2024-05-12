package router

import (
	"github.com/gin-gonic/gin"
	"online-judge/app/gateway/http"
	"online-judge/app/gateway/middleware"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors())
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		v1.POST("/submission", http.SubmissionHandler)
	}
	return ginRouter
}
