package evaluation

import "github.com/gin-gonic/gin"

type ApiEvaluation struct{}

func (e *ApiEvaluation) GetEvaluationResult(c *gin.Context) {

}

func (e *ApiEvaluation) GetUserEvaluations(c *gin.Context) {}

func (e *ApiEvaluation) GetProblemEvaluations(c *gin.Context) {}
