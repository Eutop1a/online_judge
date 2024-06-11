package user

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type ApiUser struct{}

func (u *ApiUser) InitApiUser(Router *gin.RouterGroup) {
	userApi := v1.ApiGroupApp.ApiUser

	Router.POST("/user-id", userApi.GetUserID)
	Router.POST("/detail", userApi.GetUserDetail)
	Router.POST("/update", userApi.UpdateUserDetail)
}
