package admin

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type ApiAdminProblem struct{}

func (a *ApiAdminProblem) InitAdminProblem(Router *gin.RouterGroup) {
	// 题目相关
	adminApi := v1.ApiGroupApp.ApiAdmin
	adminProblem := Router.Group("/problem")
	{
		file := adminProblem.Group("/file") // 输入输出为文件
		{
			file.POST("/create", adminApi.CreateProblemWithFile)        // 创建新题目
			file.PUT("/update", adminApi.UpdateProblemWithFile)         // 创建新题目
			file.DELETE("/:problem_id", adminApi.DeleteProblemWithFile) // 删除题目
		}
		adminProblem.POST("/create", adminApi.CreateProblem)   // 创建新题目
		adminProblem.PUT("/update", adminApi.UpdateProblem)    // 更新题目信息
		adminProblem.DELETE("/delete", adminApi.DeleteProblem) // 删除题目
	}
}
