package admin

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type ApiAdminCategory struct{}

func (c *ApiAdminCategory) InitAdminCategory(Router *gin.RouterGroup) {
	adminApi := v1.ApiGroupApp.ApiAdmin
	adminCategory := Router.Group("/category")
	{
		adminCategory.POST("/create", adminApi.AddCategory)      // 创建分类
		adminCategory.PUT("/update", adminApi.UpdateCategory)    // 更新分类信息
		adminCategory.DELETE("/delete", adminApi.DeleteCategory) // 删除分类
	}
}
