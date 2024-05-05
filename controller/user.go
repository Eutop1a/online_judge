package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online-judge/pkg"
	"online-judge/pkg/resp"
	"online-judge/services"
	"strconv"
)

// Register 用户注册接口
// @Tags User API
// @Summary 用户注册
// @Description 用户注册接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Success 200 {object} _Response "注册成功"
// @Failure 200 {object} _Response “验证码错误或已过期”
// @Failure 200 {object} _Response “该邮箱已经存在”
// @Failure 200 {object} _Response “服务器内部错误”
// @Router /register [POST]
func Register(c *gin.Context) {

	var newUser services.UserService
	//uname := c.PostForm("user_name")
	//psw := c.PostForm("password")
	//email := c.PostForm("email")
	//code := c.PostForm("code")
	//fmt.Println(uname)
	//fmt.Println(psw)
	//fmt.Println(email)
	//fmt.Println(code)
	if err := c.ShouldBind(&newUser); err != nil { //
		zap.L().Error("Register.ShouldBind error " + err.Error())
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	//
	//fmt.Println("Username", newUser.UserName)
	//fmt.Println("pwd", newUser.Password)
	//fmt.Println("email", newUser.Email)
	//fmt.Println("code", newUser.Code)

	var ret resp.RegisterResponse
	ret = newUser.Register()
	switch ret.Code {

	// 成功
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 验证码错误
	case resp.ErrorVerCode:
		resp.ResponseError(c, resp.CodeErrorVerCode)

	// 验证码过期
	case resp.ExpiredVerCode:
		resp.ResponseError(c, resp.CodeExpiredVerCode)

	// 用户名已存在
	case resp.UsernameAlreadyExist:
		resp.ResponseError(c, resp.CodeUserExist)

	// 邮箱已存在
	case resp.EmailAlreadyExist:
		resp.ResponseError(c, resp.CodeEmailExist)

	// 服务器内部错误
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// Login 用户登录接口
// @Tags User API
// @Summary 用户登录
// @Description 用户登录接口
// @Accept multipart/form-data
// @Produce json,xml
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Success 200 {object} _Response "登录成功"
// @Failure 200 {object} _Response "用户名不存在或验证码错误"
// @Failure 200 {object} _Response "验证码过期"
// @Failure 200 {object} _Response "密码错误"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /login [POST]
func Login(c *gin.Context) {
	var login services.UserService

	if err := c.ShouldBind(&login); err != nil {
		zap.L().Error("Login.ShouldBind error " + err.Error())
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	//fmt.Println("Username", login.UserName)
	//fmt.Println("pwd", login.Password)
	//fmt.Println("email", login.Email)
	//fmt.Println("code", login.Code)

	var ret resp.RegisterResponse
	ret = login.Login()
	//fmt.Println("ret.Code = ", ret.Code)
	switch ret.Code {

	// 成功，返回token
	case resp.Success:
		resp.ResponseSuccess(c, ret.Token)

	// 验证码错误
	case resp.ErrorVerCode:
		resp.ResponseError(c, resp.CodeErrorVerCode)

	// 验证码过期
	case resp.ExpiredVerCode:
		resp.ResponseError(c, resp.CodeExpiredVerCode)

	// 用户不存在
	case resp.NotExistUsername:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 密码错误
	case resp.ErrorPwd:
		resp.ResponseError(c, resp.CodeInvalidPassword)

	// 内部错误
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// GetUserDetail 获取用户详细信息接口
// @Tags User API
// @Summary 获取用户详细信息
// @Description 获取用户详细信息接口
// @Accept multipart/form-data
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} _Response "获取用户信息成功"
// @Failure 200 {object} _Response "参数错误"
// @Failure 200 {object} _Response "没有此用户ID"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /users/{user_id} [GET]
func GetUserDetail(c *gin.Context) {
	var getDetail services.UserService
	uid := c.Query("user_id")
	if uid == "" {
		zap.L().Error("GetUserDetail params error")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	getDetail.UserID, _ = strconv.ParseInt(uid, 10, 64)
	var ret resp.GetDetailResponse
	ret = getDetail.GetUserDetail()

	switch ret.Code {
	// 成功
	case resp.Success:
		resp.ResponseSuccess(c, ret.Data)

	// 用户不存在
	case resp.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 内部错误
	case resp.SearchDBError:
		resp.ResponseError(c, resp.CodeInternalServerError)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// DeleteUser 删除用户接口
// @Tags User API
// @Summary 删除用户
// @Description 删除用户接口
// @Accept multipart/form-data
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} _Response "删除用户成功"
// @Failure 200 {object} _Response "参数错误"
// @Failure 200 {object} _Response "没有此用户ID"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /users/{user_id} [DELETE]
func DeleteUser(c *gin.Context) {
	var deleteUser services.UserService
	uid := c.Query("user_id")
	if uid == "" {
		zap.L().Error("deleteUser params error")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	deleteUser.UserID, _ = strconv.ParseInt(uid, 10, 64)
	var ret resp.DeleteUserResponse
	ret = deleteUser.DeleteUser()

	switch ret.Code {
	// 成功
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 用户不存在
	case resp.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 服务器内部错误
	case resp.SearchDBError, resp.DBDeleteError:
		resp.ResponseError(c, resp.CodeInternalServerError)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// UpdateUserDetail 更新用户详细信息接口
// @Tags User API
// @Summary 更新用户详细信息
// @Description 更新用户详细信息接口
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "用户ID"
// @Param password formData string false "用户密码"
// @Param email formData string false "用户邮箱"
// @Param code formData string false "邮箱验证码"
// @Success 200 {object} _Response "更新用户信息成功"
// @Failure 200 {object} _Response "参数错误"
// @Failure 200 {object} _Response "没有此用户ID or 验证码错误"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /users/{user_id} [PUT]
func UpdateUserDetail(c *gin.Context) {
	var update services.UserService
	if err := c.ShouldBind(&update); err != nil { //
		zap.L().Error("UpdateUserDetail.ShouldBind error " + err.Error())
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	var ret resp.UpdateUserDetailResponse
	ret = update.UpdateUserDetail()

	switch ret.Code {
	// 成功
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 用户不存在
	case resp.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 验证码错误
	case resp.ErrorVerCode:
		resp.ResponseError(c, resp.CodeErrorVerCode)

	// 验证码过期
	case resp.ExpiredVerCode:
		resp.ResponseError(c, resp.CodeExpiredVerCode)

	case resp.SearchDBError, resp.EncryptPwdError:
		resp.ResponseError(c, resp.CodeInternalServerError)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// SendEmailCode 发送邮箱验证码接口
// @Tags Verification Code API
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码接口
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "邮箱"
// @Success 200 {object} _Response "发送邮箱验证码成功"
// @Failure 200 {object} _Response "邮箱格式错误"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /send-email-code [POST]
func SendEmailCode(c *gin.Context) {
	userEmail := c.PostForm("email") //从前端获取email信息
	// 判断email是否合法

	if !pkg.ValidateEmail(userEmail) {
		resp.ResponseError(c, resp.CodeInvalidateEmailFormat)
		zap.L().Error("controller-SendEmailCode-ValidateEmail " +
			fmt.Sprintf("invalid email %s ", userEmail))
		return
	}
	resCode := services.SendEmailCode(userEmail)
	switch resCode {
	// 成功
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 邮箱格式错误
	case resp.InvalidateEmailFormat:
		resp.ResponseError(c, resp.CodeInvalidateEmailFormat)

	case resp.SearchDBError, resp.EncryptPwdError:
		resp.ResponseError(c, resp.CodeInternalServerError)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// SendCode 发送图片验证码接口
// @Tags Verification Code API
// @Summary 发送图片验证码
// @Description 发送图片验证码接口
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Success 200 {object} _Response "发送图片验证码成功"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /send-code [POST]
func SendCode(c *gin.Context) {
	username := c.PostForm("username")
	b64s, err := services.SendCode(username)
	// 生成图片验证码失败
	if err != nil {
		resp.ResponseError(c, resp.CodeInternalServerError)
		return
	}
	resp.ResponseSuccess(c, b64s)

}

// CheckPictureCode 检查图片验证码接口
// @Tags Verification Code API
// @Summary 检查图片验证码
// @Description 检查图片验证码
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param code formData string true "图片验证码"
// @Success 200 {object} _Response "图片验证码正确"
// @Failure 200 {object} _Response "图片验证码错误"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /check-picture-code [POST]
func CheckPictureCode(c *gin.Context) {
	username := c.PostForm("username")
	code := c.PostForm("code")
	ok, err := services.CheckCode(username, code)
	// 从 redis 中获取失败
	if err != nil {
		resp.ResponseError(c, resp.CodeUsernameNotExist)
		return
	}
	if !ok {
		resp.ResponseError(c, resp.CodePictureError)
		return
	}
	resp.ResponseSuccess(c, resp.CodeSuccess)

}

// GetUserID 获取用户ID接口
// @Tags User API
// @Summary 获取用户ID
// @Description 获取用户ID接口
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Success 200 {object} _Response "获取用户ID成功"
// @Failure 200 {object} _Response "用户名不存在"
// @Router /user-id [POST]
func GetUserID(c *gin.Context) {
	username := c.PostForm("username")
	uid, err := services.GetUserID(username)
	if err != nil {
		resp.ResponseError(c, resp.CodeUseNotExist)
		return
	}
	resp.ResponseSuccess(c, uid)
}
