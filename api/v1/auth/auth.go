package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/models/auth/request"
	"online_judge/models/common/response"
)

type ApiAuth struct{}

// Register 用户注册接口
// @Tags Auth API
// @Summary 用户注册
// @Description 用户注册接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Success 200 {object} common.RegisterResponse "注册成功"
// @Failure 200 {object} common.RegisterResponse "用户已存在"
// @Failure 200 {object} common.RegisterResponse "验证码错误或已过期"
// @Failure 200 {object} common.RegisterResponse "验证码过期"
// @Failure 200 {object} common.RegisterResponse "该邮箱已经存在"
// @Failure 200 {object} common.RegisterResponse "服务器内部错误"
// @Router /auth/register [POST]
func (a *ApiAuth) Register(c *gin.Context) {

	var registerReq request.AuthRegisterReq
	if err := c.ShouldBind(&registerReq); err != nil { //
		zap.L().Error("Register-ShouldBind error " + err.Error())
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	//fmt.Println("Username", registerReq.UserName)
	//fmt.Println("pwd", registerReq.Password)
	//fmt.Println("email", registerReq.Email)
	//fmt.Println("code", registerReq.Code)

	var ret response.Response
	ret = AuthService.Register(registerReq)
	switch ret.Code {

	// 成功
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	// 验证码错误
	case resp_code.ErrorVerCode:
		response.ResponseError(c, response.CodeErrorVerCode)

	// 验证码过期
	case resp_code.ExpiredVerCode:
		response.ResponseError(c, response.CodeExpiredVerCode)

	// 用户名已存在
	case resp_code.UsernameAlreadyExist:
		response.ResponseError(c, response.CodeUserExist)

	// 邮箱已存在
	case resp_code.EmailAlreadyExist:
		response.ResponseError(c, response.CodeEmailExist)

	// 服务器内部错误
	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// Login 用户登录接口
// @Tags Auth API
// @Summary 用户登录
// @Description 用户登录接口
// @Accept multipart/form-data
// @Produce json,xml
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} common.LoginResponse "登录成功"
// @Failure 200 {object} common.LoginResponse "参数错误"
// @Failure 200 {object} common.LoginResponse "用户名不存在"
// @Failure 200 {object} common.LoginResponse "验证码错误"
// @Failure 200 {object} common.LoginResponse "验证码过期"
// @Failure 200 {object} common.LoginResponse "密码错误"
// @Failure 200 {object} common.LoginResponse "服务器内部错误"
// @Router /auth/login [POST]
func (a *ApiAuth) Login(c *gin.Context) {
	var loginReq request.AuthLoginReq

	if err := c.ShouldBind(&loginReq); err != nil {
		zap.L().Error("Login-ShouldBind error " + err.Error())
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	//fmt.Println("Username", loginReq.UserName)
	//fmt.Println("pwd", loginReq.Password)
	//fmt.Println("email", loginReq.Email)
	//fmt.Println("code", loginReq.Code)

	ret := AuthService.Login(loginReq)
	switch ret.Code {

	// 成功，返回token
	case resp_code.Success:
		response.ResponseSuccess(c, ret.Data)

	// 验证码错误
	case resp_code.ErrorVerCode:
		response.ResponseError(c, response.CodeErrorVerCode)

	// 验证码过期
	case resp_code.ExpiredVerCode:
		response.ResponseError(c, response.CodeExpiredVerCode)

	// 用户不存在
	case resp_code.NotExistUsername:
		response.ResponseError(c, response.CodeUseNotExist)

	// 密码错误
	case resp_code.ErrorPwd:
		response.ResponseError(c, response.CodeInvalidPassword)

	// 内部错误
	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}
