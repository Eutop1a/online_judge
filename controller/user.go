package controller

import (
	"OnlineJudge/models"
	"OnlineJudge/pkg"
	"OnlineJudge/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// Register 用户注册接口
// @Summary 用户注册
// @Description 用户注册接口
// @Accept multipart/form-data
// @Produce json
// @Param user_name formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Success 200 {object} _RegisterSuccess "注册成功"
// @Failure 400 {object} _RegisterError “验证码错误或已过期”
// @Failure 403 {object} _RegisterError “该邮箱已经存在”
// @Failure 500 {object} _RegisterError “服务器内部错误”
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "params error",
		})
		return
	}
	//
	//fmt.Println("Username", newUser.UserName)
	//fmt.Println("pwd", newUser.Password)
	//fmt.Println("email", newUser.Email)
	//fmt.Println("code", newUser.Code)

	var ret models.RegisterResponse
	ret = newUser.Register()
	switch ret.Code {
	// 200
	case services.Success:
		c.JSON(http.StatusOK, gin.H{
			"token": ret.Token,
			"msg":   "register successfully",
		})

	// 400
	case services.ErrorVerCode, services.ExpiredVerCode:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or expired verification code",
		})

		// 403
	case services.UsernameAlreadyExist:
		c.JSON(http.StatusForbidden, gin.H{
			"error": "username already exists",
		})
		// 403
	case services.EmailAlreadyExist:
		c.JSON(http.StatusForbidden, gin.H{"error": "email already exists"})
		// 500
	case services.GenerateNodeError, services.SearchDBError, services.EncryptPwdError,
		services.InsertNewUserError, services.GenerateTokenError:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "internal server error"})
	}
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户登录接口
// @Accept multipart/form-data
// @Produce json
// @Param user_name formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Success 200 {object} _LoginSuccess "登录成功"
// @Failure 400 {object} _LoginError "用户名不存在或验证码错误"
// @Failure 401 {object} _LoginError "验证码过期"
// @Failure 403 {object} _LoginError "密码错误"
// @Failure 500 {object} _LoginError "服务器内部错误"
// @Router /login [POST]
func Login(c *gin.Context) {
	var login services.UserService

	if err := c.ShouldBind(&login); err != nil { //
		zap.L().Error("Login.ShouldBind error " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "params error",
		})
		return
	}
	//fmt.Println("Username", login.UserName)
	//fmt.Println("pwd", login.Password)
	//fmt.Println("email", login.Email)
	//fmt.Println("code", login.Code)

	var ret models.RegisterResponse
	ret = login.Login()
	//fmt.Println("ret.Code = ", ret.Code)
	switch ret.Code {
	// 200
	case services.Success:
		c.JSON(http.StatusOK, gin.H{
			"token": ret.Token,
			"msg":   "login successfully",
		})
		// 400
	case services.ErrorVerCode:
		c.JSON(http.StatusBadRequest, gin.H{"error": "expired verification code"})
	case services.NotExistUsername:
		c.JSON(http.StatusBadRequest, gin.H{"error": "do not have this username"})
		// 401
	case services.ExpiredVerCode:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "expired verification code"})
		// 403
	case services.ErrorPwd:
		c.JSON(http.StatusForbidden, gin.H{"error": "error password"})
		// 500
	case services.GenerateNodeError, services.GenerateTokenError, services.SearchDBError, services.EncryptPwdError:
		c.JSON(http.StatusForbidden, gin.H{"error": "Internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

// GetUserDetail 获取用户详细信息接口
// @Summary 获取用户详细信息
// @Description 获取用户详细信息接口
// @Accept multipart/form-data
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} _GetUserDetailSuccess "获取成功"
// @Failure 400 {object} _GetUserDetailError "参数错误"
// @Failure 403 {object} _GetUserDetailError "没有此用户ID"
// @Failure 500 {object} _GetUserDetailError "服务器内部错误"
// @Router /users/{user_id} [GET]
func GetUserDetail(c *gin.Context) {
	var getDetail services.UserService
	uid := c.Query("user_id")
	if uid == "" {
		zap.L().Error("GetUserDetail params error")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "params error",
		})
		return
	}
	getDetail.UserID, _ = strconv.ParseInt(uid, 10, 64)
	var ret models.GetDetailResponse
	ret = getDetail.GetUserDetail()

	switch ret.Code {
	case services.Success:
		c.JSON(http.StatusOK, gin.H{"msg": "success", "data": ret.Data})
	case services.NotExistUserID:
		c.JSON(http.StatusForbidden, gin.H{"error": "no such userID"})
	case services.SearchDBError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// DeleteUser 删除用户接口
// @Summary 删除用户
// @Description 删除用户接口
// @Accept multipart/form-data
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} _DeleteUserSuccess "获取成功"
// @Failure 400 {object} _DeleteUserError "参数错误"
// @Failure 403 {object} _DeleteUserError "没有此用户ID"
// @Failure 500 {object} _DeleteUserError "服务器内部错误"
// @Router /users/{user_id} [DELETE]
func DeleteUser(c *gin.Context) {
	var deleteUser services.UserService
	uid := c.Query("user_id")
	if uid == "" {
		zap.L().Error("deleteUser params error")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "params error",
		})
		return
	}
	deleteUser.UserID, _ = strconv.ParseInt(uid, 10, 64)
	var ret models.DeleteUserResponse
	ret = deleteUser.DeleteUser()

	switch ret.Code {
	case services.Success:
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	case services.NotExistUserID:
		c.JSON(http.StatusForbidden, gin.H{"error": "no such userID"})
	case services.SearchDBError, services.DBDeleteError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// UpdateUserDetail 更新用户详细信息接口
// @Summary 更新用户详细信息
// @Description 更新用户详细信息接口
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "用户ID"
// @Param password formData string false "用户密码"
// @Param email formData string false "用户邮箱"
// @Success 200 {object} _UpdateUserDetailSuccess "获取成功"
// @Failure 400 {object} _UpdateUserDetailError "参数错误"
// @Failure 403 {object} _UpdateUserDetailError "没有此用户ID or 验证码错误"
// @Failure 500 {object} _UpdateUserDetailError "服务器内部错误"
// @Router /users/{user_id} [PUT]
func UpdateUserDetail(c *gin.Context) {
	var update services.UserService
	if err := c.ShouldBind(&update); err != nil { //
		zap.L().Error("UpdateUserDetail.ShouldBind error " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "params error",
		})
		return
	}
	var ret models.UpdateUserDetailResponse
	ret = update.UpdateUserDetail()

	switch ret.Code {
	case services.Success:
		c.JSON(http.StatusOK, gin.H{"msg": "Success"})
	case services.NotExistUserID:
		c.JSON(http.StatusForbidden, gin.H{"error": "No such userID"})
	case services.ExpiredVerCode:
		c.JSON(http.StatusForbidden, gin.H{"error": "expired verification code"})
	case services.ErrorVerCode:
		c.JSON(http.StatusForbidden, gin.H{"error": "error verification code"})
	case services.SearchDBError, services.EncryptPwdError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

/*
SearchDBError
NotExistUserID
ExpiredVerCode
ErrorVerCode
EncryptPwdError
SearchDBError
Success
*/

// SendCode 发送验证码接口
// @Summary 发送验证码
// @Description 发送验证码接口
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "邮箱"
// @Success 200 {object} _SendCodeSuccess "发送验证码成功"
// @Failure 400 {object} _SendCodeError "邮箱格式错误"
// @Failure 500 {object} _SendCodeError "服务器内部错误"
// @Router /sendCode [POST]
func SendCode(c *gin.Context) {
	userEmail := c.PostForm("email") //从前端获取email信息
	// 判断email是否合法
	//fmt.Println("email:", userEmail)
	if !pkg.ValidateEmail(userEmail) {
		c.String(http.StatusBadRequest, "Invalidate Email format")
		return
	}
	resCode := services.SendCode(userEmail)
	switch resCode {
	case services.Success:
		c.JSON(http.StatusOK, gin.H{"msg": "Send verification code successfully"})
	case services.InvalidateEmailFormat:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalidate email format"})
	case services.SendCodeError:
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Send verification code error"})
	}
}
