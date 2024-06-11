package submission

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type Submission struct{}

func (s *Submission) InitSubmission(Router *gin.RouterGroup) {
	submissionApi := v1.ApiGroupApp.ApiSubmission

	Router.POST("/code", submissionApi.SubmitCode)              // 提交代码
	Router.POST("/file/code", submissionApi.SubmitCodeWithFile) // 提交代码
}
