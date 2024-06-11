package problem

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type Problem struct{}

func (p *Problem) InitProblem(Router *gin.RouterGroup) {
	problemApi := v1.ApiGroupApp.ApiProblem

	Router.POST("/id", problemApi.GetProblemID)             // 获取题目ID
	Router.GET("/list", problemApi.GetProblemList)          // 获取题目列表
	Router.GET("/:problem_id", problemApi.GetProblemDetail) // 获取单个题目详细
	Router.GET("/random", problemApi.GetProblemRandom)      // 随机单个题目详细
	Router.POST("/search", problemApi.SearchProblem)        // 搜索题目
}
