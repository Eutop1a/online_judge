package services

import (
	"OnlineJudge/dao/mysql"
	"OnlineJudge/dao/redis"
	"OnlineJudge/models"
	"OnlineJudge/pkg"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"time"
)

type UserService struct {
	UserID           int64     `form:"user_id" json:"user_id"`
	UserName         string    `form:"user_name" json:"user_name" validate:"required"`
	Password         string    `form:"password" json:"password" validate:"required"`
	Email            string    `form:"email" json:"email" validate:"required"`
	Code             string    `form:"code" json:"code" validate:"required"`
	RegistrationDate time.Time `form:"registration_date" json:"registration_date"`
	LastLoginData    time.Time `form:"last_login_data" json:"last_login_data"`
	Role             bool      `form:"role" json:"role"`
	// true is Admin, false is user
}

// 用户注册服务
func (u *UserService) Register() (response models.RegisterResponse) {

	// 检查数据库中是否已经有这个email
	var countEmail, countUsername int64
	err := mysql.CheckEmail(u.Email, &countEmail)
	// UID作为唯一标识，但是email和username不能重复
	switch {
	case err != nil:
		response.Code = SearchDBError
		zap.L().Error("SearchDBError" + err.Error())
		return
	case countEmail > 0:
		response.Code = EmailAlreadyExist
		zap.L().Error(fmt.Sprintf("Email %s aleardy exist", u.Email))
		return
	}

	// 检查数据库中是否已经有这个username
	err = mysql.CheckUsername(u.UserName, &countUsername)
	switch {
	case err != nil:
		response.Code = SearchDBError
		zap.L().Error("SearchDBError" + err.Error())
		return
	case countUsername > 0:
		response.Code = UsernameAlreadyExist
		zap.L().Error(fmt.Sprintf("username %s aleardy exist", u.UserName))
		return
	}

	//// 检验邮箱是否有用
	//ok := pkg.ValidateEmail(u.Email)
	//if !ok {
	//	response.Code = InvalidateEmailFormat
	//	zap.L().Error(fmt.Sprintf("email %s invalid", u.Email))
	//	return
	//}

	//验证码获取及验证
	code, err := redis.GetVerificationCode(u.Email)
	// 验证码过期
	if errors.Is(err, fmt.Errorf("verification code expired")) {
		response.Code = ExpiredVerCode
		return
	}
	// 验证码错误
	if code != u.Code {
		response.Code = ErrorVerCode
		zap.L().Error(fmt.Sprintf("Error verfiction code %s:%s", u.Email, code))
		return
	}

	u.LastLoginData = time.Now()
	u.RegistrationDate = time.Now()

	// 生成唯一的ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		response.Code = GenerateNodeError
		zap.L().Error("generate new node error" + err.Error())
		return
	}

	u.UserID = int64(node.Generate())
	// 密码加密
	u.Password, err = pkg.CryptoPwd(u.Password)
	if err != nil {
		response.Code = EncryptPwdError
		zap.L().Error("encrypt pwd error " + err.Error())
		return
	}

	uID := u.UserID
	userName := u.UserName
	passWord := u.Password
	email := u.Email
	registrationDate := u.RegistrationDate
	lastLoginData := u.LastLoginData
	// 插入数据库
	if err = mysql.InsertNewUser(uID, userName, passWord, email, registrationDate, lastLoginData); err != nil {
		response.Code = InsertNewUserError
		zap.L().Error(fmt.Sprintf("insert new user %s error", u.UserName) + err.Error())
		return
	}
	// 生成token
	token := pkg.GenerateToken(u.UserName)
	if err != nil {
		response.Code = GenerateTokenError
		zap.L().Error("generate token error")
		return
	}
	// 设置响应的状态码和 token
	response.Code = Success
	response.Token = token
	return
}

// Login 用户登录服务
func (u *UserService) Login() (response models.RegisterResponse) {
	//验证码获取及验证
	code, err := redis.GetVerificationCode(u.Email)
	// 验证码过期
	if errors.Is(err, fmt.Errorf("verification code expired")) {
		response.Code = ExpiredVerCode
		return
	}
	// 验证码错误
	if code != u.Code {
		response.Code = ErrorVerCode
		zap.L().Error(fmt.Sprintf("Error verfiction code %s:%s", u.Email, code))
		return
	}
	// 检验是否有这个用户名
	var UserNameCount int64
	err = mysql.CheckUsername(u.UserName, &UserNameCount)
	if err != nil {
		response.Code = SearchDBError
		zap.L().Error("SearchDBError" + err.Error())
		return
	}
	if UserNameCount == 0 {
		zap.L().Error(fmt.Sprintf("Do not have this username: %s", u.UserName))
		response.Code = NotExistUsername
		return
	}

	// 检查密码是否正确
	ok, err := mysql.CheckPwd(u.UserName, u.Password)
	if err != nil {
		response.Code = SearchDBError
		zap.L().Error("Search DB error: " + err.Error())
		return
	}
	// 密码错误
	if !ok {
		response.Code = ErrorPwd
		zap.L().Error("Password error")
		return
	}
	u.LastLoginData = time.Now()
	T, err := mysql.UpdateLoginData(u.UserName, u.LastLoginData)
	if err != nil {
		if T == 0 {
			response.Code = SearchDBError
			zap.L().Error("Find user error: " + err.Error())
			return
		} else if T == -1 {
			response.Code = DBSaveError
			zap.L().Error("Save to DB: " + err.Error())
			return
		}
	}
	// 成功返回
	// 生成token
	token := pkg.GenerateToken(u.UserName)
	if err != nil {
		response.Code = GenerateTokenError
		zap.L().Error("generate token error")
		return
	}
	// 设置响应的状态码和 token
	response.Code = Success
	response.Token = token
	return
}

// 获取用户详细信息
func (u *UserService) GetUserDetail() {

}

// 删除用户
func (u *UserService) DeleteUser() {

}

// 更新用户详细信息
func (u *UserService) UpdateUserDetail() {

}

// 发送验证码接口
func SendCode(userEmail string) (resCode int) {
	// 验证邮箱是否合法
	if !pkg.ValidateEmail(userEmail) {
		resCode = InvalidateEmailFormat
		return
	}
	// 创建验证码
	code, ts := pkg.CreateVerificationCode()
	// 发送验证码
	err := pkg.SendCode(userEmail, code)
	if err != nil {
		resCode = SendCodeError
		return
	}
	// redis保存email和验证码的键值对
	err = redis.StoreVerificationCode(userEmail, code, ts)
	if err != nil {
		resCode = StoreVerCodeError
		return
	}

	resCode = Success
	return
}
