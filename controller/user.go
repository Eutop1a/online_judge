package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online-judge/consts/resp_code"
	"online-judge/pkg/resp"
	"online-judge/services"
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
// @Success 200 {object} models.RegisterResponse "注册成功"
// @Failure 200 {object} models.RegisterResponse "用户已存在"
// @Failure 200 {object} models.RegisterResponse "验证码错误或已过期"
// @Failure 200 {object} models.RegisterResponse "验证码过期"
// @Failure 200 {object} models.RegisterResponse "该邮箱已经存在"
// @Failure 200 {object} models.RegisterResponse "服务器内部错误"
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

	var ret resp.Response
	ret = newUser.Register()
	switch ret.Code {

	// 成功
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 验证码错误
	case resp_code.ErrorVerCode:
		resp.ResponseError(c, resp.CodeErrorVerCode)

	// 验证码过期
	case resp_code.ExpiredVerCode:
		resp.ResponseError(c, resp.CodeExpiredVerCode)

	// 用户名已存在
	case resp_code.UsernameAlreadyExist:
		resp.ResponseError(c, resp.CodeUserExist)

	// 邮箱已存在
	case resp_code.EmailAlreadyExist:
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
// @Success 200 {object} models.LoginResponse "登录成功"
// @Failure 200 {object} models.LoginResponse "参数错误"
// @Failure 200 {object} models.LoginResponse "用户名不存在"
// @Failure 200 {object} models.LoginResponse "验证码错误"
// @Failure 200 {object} models.LoginResponse "验证码过期"
// @Failure 200 {object} models.LoginResponse "密码错误"
// @Failure 200 {object} models.LoginResponse "服务器内部错误"
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

	var ret resp.ResponseWithData
	ret = login.Login()
	//fmt.Println("ret.Code = ", ret.Code)
	switch ret.Code {

	// 成功，返回token
	case resp_code.Success:
		resp.ResponseSuccess(c, ret.Data)

	// 验证码错误
	case resp_code.ErrorVerCode:
		resp.ResponseError(c, resp.CodeErrorVerCode)

	// 验证码过期
	case resp_code.ExpiredVerCode:
		resp.ResponseError(c, resp.CodeExpiredVerCode)

	// 用户不存在
	case resp_code.NotExistUsername:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 密码错误
	case resp_code.ErrorPwd:
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
// // @Param user_id path string true "用户ID"
// @Param Authorization header string true "token"
// @Success 200 {object} models.GetUserDetailResponse "获取用户信息成功"
// @Failure 200 {object} models.GetUserDetailResponse "参数错误"
// @Failure 200 {object} models.GetUserDetailResponse "没有此用户ID"
// @Failure 200 {object} models.GetUserDetailResponse "服务器内部错误"
// @Router /users/get [GET]
func GetUserDetail(c *gin.Context) {
	var getDetail services.UserService
	uid, ok := c.Get(resp.CtxUserIDKey)
	if !ok {
		resp.ResponseError(c, resp.CodeNeedLogin)
		return
	}
	//uid := c.Param("user_id")
	//if uid == "" {
	//	zap.L().Error("GetUserDetail params error")
	//	resp.ResponseError(c, resp.CodeInvalidParam)
	//	return
	//}
	//getDetail.UserID, _ = strconv.ParseInt(uid, 10, 64)
	getDetail.UserID = uid.(int64)
	var ret resp.ResponseWithData
	ret = getDetail.GetUserDetail()

	switch ret.Code {
	// 成功
	case resp_code.Success:
		resp.ResponseSuccess(c, ret.Data)

	// 用户不存在
	case resp_code.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 内部错误
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
// @Param Authorization header string true "token"
// @Param username formData string false "用户名"
// @Param password formData string false "用户密码"
// @Param email formData string false "用户邮箱"
// @Param code formData string false "邮箱验证码"
// @Success 200 {object} models.UpdateUserDetailResponse "更新用户信息成功"
// @Failure 200 {object} models.UpdateUserDetailResponse "参数错误"
// @Failure 200 {object} models.UpdateUserDetailResponse "没有此用户ID"
// @Failure 200 {object} models.UpdateUserDetailResponse "验证码错误"
// @Failure 200 {object} models.UpdateUserDetailResponse "验证码过期"
// @Failure 200 {object} models.UpdateUserDetailResponse "服务器内部错误"
// @Router /users/update [PUT]
func UpdateUserDetail(c *gin.Context) {
	var update services.UserService
	if err := c.ShouldBind(&update); err != nil { //
		zap.L().Error("UpdateUserDetail.ShouldBind error " + err.Error())
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	uid, ok := c.Get(resp.CtxUserIDKey)
	if !ok {
		resp.ResponseError(c, resp.CodeNeedLogin)
		return
	}
	//uid := c.Param("user_id")
	//if uid == "" {
	//	zap.L().Error("UpdateUserDetail params error")
	//	resp.ResponseError(c, resp.CodeInvalidParam)
	//	return
	//}
	//update.UserID, _ = strconv.ParseInt(uid, 10, 64)
	//fmt.Println(update.UserID)
	update.UserID = uid.(int64)
	var ret resp.Response
	ret = update.UpdateUserDetail()

	switch ret.Code {
	// 成功
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 用户不存在
	case resp_code.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseNotExist)

	// 验证码错误
	case resp_code.ErrorVerCode:
		resp.ResponseError(c, resp.CodeErrorVerCode)

	// 验证码过期
	case resp_code.ExpiredVerCode:
		resp.ResponseError(c, resp.CodeExpiredVerCode)

	// 新用户名已经存在
	case resp_code.UsernameAlreadyExist:
		resp.ResponseError(c, resp.CodeUsernameAlreadyExist)

	// 未申请验证码
	case resp_code.NeedObtainVerificationCode:
		resp.ResponseError(c, resp.CodeObtainVerificationCode)

	// 邮箱已经存在
	case resp_code.EmailAlreadyExist:
		resp.ResponseError(c, resp.CodeEmailExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// GetUserID 获取用户ID接口
// @Tags User API
// @Summary 获取用户ID
// @Description 获取用户ID接口
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Success 200 {object} models.GetUserIDResponse "获取用户ID成功"
// @Failure 200 {object} models.GetUserIDResponse "用户名不存在"
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
