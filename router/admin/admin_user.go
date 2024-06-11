package admin

import (
	"github.com/gin-gonic/gin"
	"online_judge/api/v1"
)

type ApiAdminUser struct{}

// InitAdminUser 用户相关
func (a *ApiAdminUser) InitAdminUser(Router *gin.RouterGroup) {
	adminApi := v1.ApiGroupApp.ApiAdmin
	adminUsers := Router.Group("/users")
	{
		adminUsers.DELETE("/:user_id", adminApi.DeleteUser) // 删除用户
		adminUsers.POST("/add-admin", adminApi.AddAdmin)    // 添加用户为管理员
	}
}
