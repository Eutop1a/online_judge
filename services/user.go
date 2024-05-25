package services

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"online-judge/consts"
	"online-judge/dao/mysql"
	"online-judge/dao/redis"
	"online-judge/pkg/jwt"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
)

type UserService struct {
	UserID   int64  `form:"user_id" json:"user_id"`
	UserName string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Email    string `form:"email" json:"email" validate:"required"`
	Code     string `form:"code" json:"code" validate:"required"`
	//Role             bool      `form:"role" json:"role"`
	// true is Admin, false is user
}

// Register 用户注册服务
func (u *UserService) Register() (response resp.Response) {

	// 检查数据库中是否已经有这个email
	var countEmail, countUsername int64
	err := mysql.CheckEmail(u.Email, &countEmail)
	// UID作为唯一标识，但是email和username不能重复
	switch {
	case err != nil:
		response.Code = consts.SearchDBError
		zap.L().Error("services-Register-CheckEmail searchDBError", zap.Error(err))
		return
	case countEmail > 0:
		response.Code = consts.EmailAlreadyExist
		zap.L().Error("services-Register-CheckEmail " +
			fmt.Sprintf("email %s aleardy exist ", u.Email))
		return
	}

	// 检查数据库中是否已经有这个username
	err = mysql.CheckUsername(u.UserName, &countUsername)
	switch {
	case err != nil:
		response.Code = consts.SearchDBError
		zap.L().Error("services-Register-CheckUsername ", zap.Error(err))
		return
	case countUsername > 0:
		response.Code = consts.UsernameAlreadyExist
		zap.L().Error("services-Register-CheckUsername " +
			fmt.Sprintf("username %s aleardy exist ", u.UserName))
		return
	}

	//// 检验邮箱是否有用
	//ok := pkg.ValidateEmail(u.Email)
	//if !ok {
	//	resp.Code = InvalidateEmailFormat
	//	zap.L().Error(fmt.Sprintf("email %s invalid", u.Email))
	//	return
	//}

	//验证码获取及验证
	code, err := redis.GetVerificationCode(u.Email)
	// 验证码过期
	if errors.Is(err, fmt.Errorf("verification code expired")) {
		response.Code = consts.ExpiredVerCode
		zap.L().Error("services-Register-GetVerificationCode " +
			fmt.Sprintf("verfiction code expired %s, %s ", u.Email, code))
		return
	}
	// 验证码错误
	if code != u.Code {
		response.Code = consts.ErrorVerCode
		zap.L().Error("services-Register-GetVerificationCode " +
			fmt.Sprintf("error verfiction code %s:%s ", u.Email, code))
		return
	}

	// 生成唯一的ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		response.Code = consts.GenerateNodeError
		zap.L().Error("services-Register-NewNode ", zap.Error(err))
		return
	}

	u.UserID = int64(node.Generate())
	// 密码加密
	u.Password, err = utils.CryptoPwd(u.Password)
	if err != nil {
		response.Code = consts.EncryptPwdError
		zap.L().Error("services-Register-CryptoPwd ", zap.Error(err))
		return
	}

	uID := u.UserID
	userName := u.UserName
	passWord := u.Password
	email := u.Email
	// 插入数据库
	if err = mysql.InsertNewUser(uID, userName, passWord, email); err != nil {
		response.Code = consts.InsertNewUserError
		zap.L().Error("services-Register-InsertNewUser "+
			fmt.Sprintf("insert new user %s error ", u.UserName), zap.Error(err))
		return
	}
	//// 生成token
	//token, err := jwt.GenerateToken(u.UserID, u.UserName)
	//if err != nil {
	//	response.Code = resp.GenerateTokenError
	//	zap.L().Error("generate token error")
	//	return
	//}
	// 设置响应的状态码
	response.Code = consts.Success
	//response.Token = token
	return
}

// Login 用户登录服务
func (u *UserService) Login() (response resp.ResponseWithData) {

	////验证码获取及验证
	//code, err := redis.GetVerificationCode(u.Email)
	//// 验证码过期
	//if errors.Is(err, fmt.Errorf("verification code expired")) {
	//	response.Code = resp.ExpiredVerCode
	//	zap.L().Error("services-Login-GetVerificationCode " +
	//		fmt.Sprintf("verfiction code expired %s, %s ", u.Email, code))
	//	return
	//}
	//// 验证码错误
	//if code != u.Code {
	//	response.Code = resp.ErrorVerCode
	//	zap.L().Error("services-Login-GetVerificationCode " +
	//		fmt.Sprintf("error verfiction code %s:%s ", u.Email, code))
	//	return
	//}
	// 检验是否有这个用户名
	var UserNameCount int64
	err := mysql.CheckUsername(u.UserName, &UserNameCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-Login-CheckUsername ", zap.Error(err))
		return
	}
	if UserNameCount == 0 {
		response.Code = consts.NotExistUsername
		zap.L().Error("services-Login-CheckUsername " +
			fmt.Sprintf("do not have this username: %s ", u.UserName))
		return
	}

	// 检查密码是否正确
	err = mysql.CheckPwd(u.UserName, u.Password)
	if err != nil {
		response.Code = consts.ErrorPwd
		zap.L().Error("services-Login-CheckPwd ", zap.Error(err))
		return
	}
	// 成功返回
	// 生成token
	u.UserID, err = mysql.GetUserID(u.UserName)

	var userIsAdmin bool
	err = mysql.CheckUserIsAdmin(u.UserID)
	if err != nil {
		userIsAdmin = false
	}
	userIsAdmin = true
	token, err := jwt.GenerateToken(u.UserID, u.UserName, userIsAdmin)
	if err != nil {
		response.Code = consts.GenerateTokenError
		zap.L().Error("service-Login-GenerateToken: ", zap.Error(err))
		return
	}
	// 设置响应的状态码和 token
	response.Code = consts.Success
	response.Data = token
	return
}

// GetUserDetail 获取用户详细信息
func (u *UserService) GetUserDetail() (response resp.ResponseWithData) {
	// 检验是否有这个用户ID
	var UserIDCount int64
	err := mysql.CheckUserID(u.UserID, &UserIDCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-GetUserDetail-CheckUserID ", zap.Error(err))
		return
	}
	if UserIDCount == 0 {
		response.Code = consts.NotExistUserID
		zap.L().Error("services-GetUserDetail-CheckUserID " +
			fmt.Sprintf("do not have this userID: %d ", u.UserID))
		return
	}
	// 去数据库获取详细信息
	data, err := mysql.GetUserDetail(u.UserID)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-GetUserDetail-GetUserDetail "+
			fmt.Sprintf("db scan userID %d detail information failed ", u.UserID), zap.Error(err))
		return
	}
	response.Data = data
	response.Code = consts.Success
	return
}

// UpdateUserDetail 更新用户详细信息
func (u *UserService) UpdateUserDetail() (response resp.Response) {
	// 检验是否有这个用户ID
	var UserIDCount int64
	err := mysql.CheckUserID(u.UserID, &UserIDCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-UpdateUserDetail-CheckUserID ", zap.Error(err))
		return
	}
	if UserIDCount == 0 {
		response.Code = consts.NotExistUserID
		zap.L().Error("services-UpdateUserDetail-CheckUserID " +
			fmt.Sprintf("do not have this userID %d ", u.UserID))
		return
	}
	// 更新用户信息
	if u.Email != "" {
		//验证码获取及验证
		code, err := redis.GetVerificationCode(u.Email)
		// 验证码过期
		if errors.Is(err, fmt.Errorf("verification code expired")) {
			response.Code = consts.ExpiredVerCode
			zap.L().Error("services-UpdateUserDetail-GetVerificationCode "+
				fmt.Sprintf("expired verfiction code %s:%s ", u.Email, code), zap.Error(err))
			return
		}
		// 验证码错误
		if code != u.Code {
			response.Code = consts.ErrorVerCode
			zap.L().Error("services-UpdateUserDetail-GetVerificationCode "+
				fmt.Sprintf("error verfiction code %s:%s ", u.Email, code), zap.Error(err))
			return
		}
	}

	if u.Password != "" {
		// 密码加密
		u.Password, err = utils.CryptoPwd(u.Password)
		if err != nil {
			response.Code = consts.EncryptPwdError
			zap.L().Error("services-UpdateUserDetail-CryptoPwd ", zap.Error(err))
			return
		}
	}
	err = mysql.UpdateUserDetail(u.UserID, u.Email, u.Password)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-UpdateUserDetail-UpdateUserDetail "+
			fmt.Sprintf("db update userID %d failed ", u.UserID), zap.Error(err))
		return
	}
	response.Code = consts.Success
	return
}

func GetUserID(username string) (uid int64, err error) {
	uid, err = mysql.GetUserID(username)
	if err != nil {
		zap.L().Error("services-GetUserID-GetUserID "+username+"err: ", zap.Error(err))
		return 0, err
	}
	return uid, nil
}
