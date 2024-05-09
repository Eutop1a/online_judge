package controller

import (
	"github.com/gin-gonic/gin"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	"online-judge/services"
	"strconv"
	"time"
)

// SubmitCode 提交代码接口
// @Tags Submission API
// @Summary 提交代码
// @Description 提交代码接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param user_id formData string true "用户id"
// @Param problem_id formData string true "题目id"
// @Param language formData string true "语言"
// @Param code formData string true "代码"
// @Success 200 {object} _Response "提交代码成功"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /submissions/code [POST]
func SubmitCode(c *gin.Context) {
	var submission services.Submission
	userId, _ := strconv.Atoi(c.PostForm("user_id"))
	submission.ProblemID = c.PostForm("problem_id")
	submission.Language = c.PostForm("language")
	submission.Code = c.PostForm("code")

	submission.SubmissionID = utils.GetUUID()
	submission.UserID = int64(userId)
	submission.SubmissionTime = time.Now()
	response := submission.SubmitCode()
	switch response.Code {
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)
	case resp.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseIDNotExist)
	case resp.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
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
