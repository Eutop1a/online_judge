package problem

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type Problem struct{}

func (p *Problem) InitProblem(Router *gin.RouterGroup) {
	problemApi := v1.ApiGroupApp.ApiProblem

	Router.POST("/id", problemApi.GetProblemID)                          // 获取题目ID
	Router.GET("/list", problemApi.GetProblemList)                       // 获取题目列表
	Router.GET("/:problem_id", problemApi.GetProblemDetail)              // 获取单个题目详细
	Router.GET("/random", problemApi.GetProblemRandom)                   // 随机单个题目详细
	Router.POST("/title/search", problemApi.SearchProblem)               // 搜索题目
	Router.POST("/category/search", problemApi.GetProblemListByCategory) // 根据题目分类搜索题目
	Router.GET("/category-list", problemApi.GetCategoryList)             // 获取分类列表
	Router.GET("/hot-search", problemApi.GetHotSearches)                 // 获取最热搜索
	Router.GET("/recent-search", problemApi.GetRecentSearches)           // 获取最近搜索
}
