package routes

import "github.com/gin-gonic/gin"

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	return r
}
