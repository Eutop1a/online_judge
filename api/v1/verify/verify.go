package verify

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/models/common/response"
	"online_judge/models/verify/request"
	"online_judge/pkg/utils"
)

type ApiVerify struct{}

// SendEmailCode 发送邮箱验证码接口
// @Tags Verification Code API
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码接口
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "邮箱"
// @Success 200 {object} common.SendEmailCodeResponse "发送邮箱验证码成功"
// @Failure 200 {object} common.SendEmailCodeResponse "邮箱格式错误"
// @Failure 200 {object} common.SendEmailCodeResponse "服务器内部错误"
// @Router /send-email-code [POST]
func (v *ApiVerify) SendEmailCode(c *gin.Context) {
	var req request.SendEmailCodeReq
	req.UserEmail = c.PostForm("email") //从前端获取email信息

	// 判断email是否合法
	if !utils.ValidateEmail(req.UserEmail) {
		response.ResponseError(c, response.CodeInvalidateEmailFormat)
		zap.L().Error("controller-SendEmailCode-ValidateEmail " +
			fmt.Sprintf("invalid email %s ", req.UserEmail))
		return
	}
	resCode := VerifyService.SendEmailCode(req)
	switch resCode {
	// 成功
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// SendPictureCode 发送图片验证码接口
// @Tags Verification Code API
// @Summary 发送图片验证码
// @Description 发送图片验证码接口
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Success 200 {object} common.SendPictureCodeResponse "1000 发送图片验证码成功"
// @Failure 200 {object} common.SendPictureCodeResponse "1014 服务器内部错误"
// @Router /send-picture-code [POST]
func (v *ApiVerify) SendPictureCode(c *gin.Context) {
	var req request.SendPictureCodeReq
	req.Username = c.PostForm("username")
	if req.Username == "" {
		response.ResponseError(c, response.CodeNeedUsername)
		return
	}
	b64s, err := VerifyService.SendPictureCode(req)
	// 生成图片验证码失败
	if err != nil {
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	response.ResponseSuccess(c, b64s)

}

// CheckPictureCode 检查图片验证码接口
// @Tags Verification Code API
// @Summary 检查图片验证码
// @Description 检查图片验证码
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param code formData string true "图片验证码"
// @Success 200 {object} common.CheckPictureCodeResponse "图片验证码正确"
// @Failure 200 {object} common.CheckPictureCodeResponse "图片验证码错误"
// @Failure 200 {object} common.CheckPictureCodeResponse "用户名不存在"
// @Router /check-picture-code [POST]
func (v *ApiVerify) CheckPictureCode(c *gin.Context) {
	var req request.CheckCodeReq
	req.Username = c.PostForm("username")
	req.Code = c.PostForm("code")

	if req.Username == "" || req.Code == "" {
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	ok, err := VerifyService.CheckCode(req)
	// 从 redis 中获取失败
	if err != nil {
		response.ResponseError(c, response.CodeUsernameNotExist)
		return
	}
	if !ok {
		response.ResponseError(c, response.CodePictureError)
		return
	}
	response.ResponseSuccess(c, response.CodeSuccess)

}
