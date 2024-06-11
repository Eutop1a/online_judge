package evaluation

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type Evaluation struct{}

func (e *Evaluation) InitEvaluate(RouterGroup *gin.RouterGroup) {
	evaluationApi := v1.ApiGroupApp.ApiEvaluation

	RouterGroup.GET("/:id", evaluationApi.GetEvaluationResult)                   // 获取评测结果
	RouterGroup.GET("/user/:user_id", evaluationApi.GetUserEvaluations)          // 获取用户的评测记录
	RouterGroup.GET("/problem/:problem_id", evaluationApi.GetProblemEvaluations) // 获取题目的评测记录
}
