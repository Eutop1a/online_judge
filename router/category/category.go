package category

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type Category struct{}

// InitCategory 分类相关
func (c *Category) InitCategory(RouterGroup *gin.RouterGroup) {
	categoryApi := v1.ApiGroupApp.ApiCategory

	RouterGroup.GET("/get-by-category", categoryApi.GetProblemByCategory) // 获取分类下的所有题目

}
