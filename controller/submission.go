package controller

import (
	"github.com/gin-gonic/gin"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	"online-judge/services"
	"time"
)

// SubmitCode 提交代码接口
// @Tags Submission API
// @Summary 提交代码
// @Description 提交代码接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id formData string true "题目id"
// @Param language formData string true "语言"
// @Param code formData string true "代码"
// @Success 200 {object} models.SubmitCodeResponse "提交代码成功"
// @Failure 200 {object} models.SubmitCodeResponse "用户ID不存在"
// @Failure 200 {object} models.SubmitCodeResponse "题目ID不存在"
// @Failure 200 {object} models.SubmitCodeResponse "不支持的语言类型"
// @Failure 200 {object} models.SubmitCodeResponse "需要登录"
// @Failure 200 {object} models.SubmitCodeResponse "服务器内部错误"
// @Router /submissions/code [POST]
func SubmitCode(c *gin.Context) {
	var submission services.Submission
	userId, ok := c.Get(resp.CtxUserIDKey)
	if !ok {
		resp.ResponseError(c, resp.CodeNeedLogin)
		return
	}

	submission.ProblemID = c.PostForm("problem_id")
	submission.Language = c.PostForm("language")
	submission.Code = c.PostForm("code")

	submission.SubmissionID = utils.GetUUID()
	submission.UserID = userId.(int64)
	submission.SubmissionTime = time.Now()
	response := submission.SubmitCode()
	switch response.Code {
	case resp.Success:
		resp.ResponseSuccess(c, response.Data)

	case resp.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseIDNotExist)

	case resp.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)

	case resp.UnsupportedLanguage:
		resp.ResponseError(c, resp.CodeUnsupportedLanguage)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

func GetSubmissionDetail(c *gin.Context) {

}

func GetUserSubmissions(c *gin.Context) {

}

func GetProblemSubmissions(c *gin.Context) {

}
