package auth

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type ApiAuth struct{}

func (a *ApiAuth) InitAuth(routerGroup *gin.RouterGroup) {
	authApi := v1.ApiGroupApp.ApiAuth

	routerGroup.POST("/login", authApi.Login)
	routerGroup.POST("/register", authApi.Register)
}
