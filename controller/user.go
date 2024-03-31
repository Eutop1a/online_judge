package controller

import (
	"OnlineJudge/models"
	"OnlineJudge/pkg"
	"OnlineJudge/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
// @Success 200 {object} _ResponseRegister "注册成功"
// @Failure 400 {object} _ResponseRegister “验证码错误或已过期”
// @Failure 403 {object} _ResponseRegister “该邮箱已经存在”
// @Failure 500 {object} _ResponseRegister “服务器内部错误”
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
		zap.L().Error("ShouldBind error " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "params error",
		})
		return
	}

	fmt.Println("Username", newUser.UserName)
	fmt.Println("pwd", newUser.Password)
	fmt.Println("email", newUser.Email)
	fmt.Println("code", newUser.Code)

	var ret models.RegisterResponse
	ret = newUser.Register()
	switch ret.Code {
	// 200
	case services.Success:
		c.JSON(http.StatusOK, gin.H{
			"token": ret.Token,
			"msg":   "Register successfully",
		})

	// 400
	case services.ErrorVerCode, services.ExpiredVerCode:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired verification code",
		})

		// 403
	case services.UsernameAlreadyExist:
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Username already exists",
		})
		// 403
	case services.EmailAlreadyExist:
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Email already exists",
		})
		// 500
	case services.GenerateNodeError, services.SearchDBError, services.EncryptPwdError,
		services.InsertNewUserError, services.GenerateTokenError:
		//case services.SendCodeError:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "InternalServerError"})
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
// @Success 200 {object} _ResponseRegister "登录成功"
// @Failure 400 {object} _ResponseRegister "用户名不存在或验证码错误"
// @Failure 401 {object} _ResponseRegister "验证码过期"
// @Failure 403 {object} _ResponseRegister "密码错误"
// @Failure 500 {object} _ResponseRegister "服务器内部错误"
// @Router /login [POST]
func Login(c *gin.Context) {
	var login services.UserService

	if err := c.ShouldBind(&login); err != nil { //
		zap.L().Error("ShouldBind error " + err.Error())
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
			"msg":   "Login successfully",
		})
		// 400
	case services.ErrorVerCode, services.NotExistUsername:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expired verification code"})
		// 401
	case services.ExpiredVerCode:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired verification code"})
		// 403
	case services.ErrorPwd:
		c.JSON(http.StatusForbidden, gin.H{"error": "Error password"})
		// 500
	case services.GenerateNodeError, services.GenerateTokenError, services.SearchDBError, services.EncryptPwdError:
		c.JSON(http.StatusForbidden, gin.H{"error": "Internal server error"})
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Internal server error"})
	}
}

// 获取用户详细信息接口
func GetUserDetail(c *gin.Context) {

}

// 删除用户接口
func DeleteUser(c *gin.Context) {

}

// 更新用户详细信息接口
func UpdateUserDetail(c *gin.Context) {

}

// SendCode 发送验证码接口
// @Summary 发送验证码
// @Description 发送验证码接口
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "邮箱"
// @Success 200 {object} _ResponseSendCode "发送验证码成功"
// @Failure 400 {object} _ResponseSendCode "邮箱格式错误"
// @Failure 500 {object} _ResponseSendCode "服务器内部错误"
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
