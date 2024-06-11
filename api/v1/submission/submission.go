package submission

import (
	"github.com/gin-gonic/gin"
	"online_judge/consts/resp_code"
	"online_judge/models/common/response"
	"online_judge/models/submission/request"
	"online_judge/pkg/utils"
	"time"
)

type ApiSubmission struct{}

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
// @Success 200 {object} common.SubmitCodeResponse "提交代码成功"
// @Failure 200 {object} common.SubmitCodeResponse "用户ID不存在"
// @Failure 200 {object} common.SubmitCodeResponse "题目ID不存在"
// @Failure 200 {object} common.SubmitCodeResponse "不支持的语言类型"
// @Failure 200 {object} common.SubmitCodeResponse "需要登录"
// @Failure 200 {object} common.SubmitCodeResponse "服务器内部错误"
// @Router /submission/code [POST]
func (s *ApiSubmission) SubmitCode(c *gin.Context) {
	var submissionReq request.SubmissionReq
	userId, ok := c.Get(response.CtxUserIDKey)
	if !ok {
		response.ResponseError(c, response.CodeNeedLogin)
		return
	}

	submissionReq.ProblemID = c.PostForm("problem_id")
	submissionReq.Language = c.PostForm("language")
	submissionReq.Code = c.PostForm("code")

	submissionReq.SubmissionID = utils.GetUUID()
	submissionReq.UserID = userId.(int64)
	submissionReq.SubmissionTime = time.Now()

	resp := SubmissionService.SubmitCode(submissionReq)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, resp.Data)

	case resp_code.NotExistUserID:
		response.ResponseError(c, response.CodeUseIDNotExist)

	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)

	case resp_code.UnsupportedLanguage:
		response.ResponseError(c, response.CodeUnsupportedLanguage)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// SubmitCodeWithFile 提交代码接口(文件)
// @Tags Submission API
// @Summary 提交代码(文件)
// @Description 提交代码接口(文件)
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id formData string true "题目id"
// @Param language formData string true "语言"
// @Param code formData string true "代码"
// @Success 200 {object} common.SubmitCodeResponse "提交代码成功"
// @Failure 200 {object} common.SubmitCodeResponse "用户ID不存在"
// @Failure 200 {object} common.SubmitCodeResponse "题目ID不存在"
// @Failure 200 {object} common.SubmitCodeResponse "不支持的语言类型"
// @Failure 200 {object} common.SubmitCodeResponse "需要登录"
// @Failure 200 {object} common.SubmitCodeResponse "服务器内部错误"
// @Router /submission/file/code [POST]
func (s *ApiSubmission) SubmitCodeWithFile(c *gin.Context) {
	var req request.SubmissionReq
	userId, ok := c.Get(response.CtxUserIDKey)
	if !ok {
		response.ResponseError(c, response.CodeNeedLogin)
		return
	}

	req.ProblemID = c.PostForm("problem_id")
	req.Language = c.PostForm("language")
	req.Code = c.PostForm("code")

	req.SubmissionID = utils.GetUUID()
	req.UserID = userId.(int64)
	req.SubmissionTime = time.Now()
	resp := SubmissionService.SubmitCodeWithFile(req)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, resp.Data)

	case resp_code.NotExistUserID:
		response.ResponseError(c, response.CodeUseIDNotExist)

	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)

	case resp_code.UnsupportedLanguage:
		response.ResponseError(c, response.CodeUnsupportedLanguage)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

func (s *ApiSubmission) GetSubmissionDetail(c *gin.Context) {

}

func (s *ApiSubmission) GetUserSubmissions(c *gin.Context) {

}

func (s *ApiSubmission) GetProblemSubmissions(c *gin.Context) {

}
