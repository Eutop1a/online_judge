package services

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"online-judge/consts/resp_code"
	"online-judge/dao/mysql"
	"online-judge/dao/redis"
	"online-judge/pkg/jwt"
	"online-judge/pkg/resp"
	"online-judge/pkg/snowflake"
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
	// 检查数据库中是否已经有这个email或username
	var countEmail, countUsername int64
	err := mysql.CheckEmailAndUsername(u.Email, u.UserName, &countEmail, &countUsername)
	// UID作为唯一标识，但是email和username不能重复
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-Register-CheckEmailAndUsername ", zap.Error(err))
		return
	}

	if countEmail > 0 {
		response.Code = resp_code.EmailAlreadyExist
		zap.L().Warn("services-Register-CheckEmailAndUsername email already exist", zap.String("email", u.Email))
		return
	}

	if countUsername > 0 {
		response.Code = resp_code.UsernameAlreadyExist
		zap.L().Warn("services-Register-CheckEmailAndUsername username already exist", zap.String("username", u.UserName))
		return
	}

	//验证码获取及验证
	code, err := redis.GetVerificationCode(u.Email)
	// 验证码过期
	if err != nil {
		if errors.Is(err, fmt.Errorf("verification code expired")) {
			response.Code = resp_code.ExpiredVerCode
			zap.L().Warn("services-Register-GetVerificationCode verification code expired",
				zap.String("email", u.Email))
		} else {
			response.Code = resp_code.ErrorVerCode
			zap.L().Warn("services-Register-GetVerificationCode error verification code",
				zap.String("email", u.Email))
		}
		return
	}
	// 验证码错误
	if code != u.Code {
		response.Code = resp_code.ErrorVerCode
		zap.L().Warn("services-Register-GetVerificationCode incorrect verification code",
			zap.String("email", u.Email))
		return
	}

	//
	//node, err := snowflake.GetID()
	////node, err := snowflake.NewNode(1)
	//if err != nil {
	//	response.Code = resp_code.GenerateNodeError
	//	zap.L().Error("services-Register-NewNode ", zap.Error(err))
	//	return
	//}

	// 生成唯一的ID
	Id := snowflake.GetID()
	u.UserID = Id
	// 密码加密
	u.Password, err = utils.CryptoPwd(u.Password)
	if err != nil {
		response.Code = resp_code.EncryptPwdError
		zap.L().Error("services-Register-CryptoPwd ", zap.Error(err))
		return
	}

	// 插入数据库
	if err = mysql.InsertNewUser(u.UserID, u.UserName, u.Password, u.Email); err != nil {
		response.Code = resp_code.InsertNewUserError
		zap.L().Error("services-Register-InsertNewUser "+
			fmt.Sprintf("insert new user %s error ", u.UserName), zap.Error(err))
		return
	}
	// 设置响应的状态码
	response.Code = resp_code.Success
	return
}

// Login 用户登录服务
func (u *UserService) Login() (response resp.ResponseWithData) {
	//// 检验是否有这个用户名
	//var UserNameCount int64
	//err := mysql.CheckUsername(u.UserName, &UserNameCount)
	//if err != nil {
	//	response.Code = consts.SearchDBError
	//	zap.L().Error("services-Login-CheckUsername ", zap.Error(err))
	//	return
	//}
	//if UserNameCount == 0 {
	//	response.Code = consts.NotExistUsername
	//	zap.L().Error("services-Login-CheckUsername " +
	//		fmt.Sprintf("do not have this username: %s ", u.UserName))
	//	return
	//}
	//
	//// 检查密码是否正确
	//err = mysql.CheckPwd(u.UserName, u.Password)
	//if err != nil {
	//	response.Code = consts.ErrorPwd
	//	zap.L().Error("services-Login-CheckPwd ", zap.Error(err))
	//	return
	//}
	//// 成功返回
	//// 生成token
	//u.UserID, err = mysql.GetUserID(u.UserName)
	//
	//var userIsAdmin bool
	//err = mysql.CheckUserIsAdminByUsername(u.UserName)
	//if err != nil {
	//	userIsAdmin = false
	//}
	//userIsAdmin = true

	// 检查用户名和密码是否正确
	userID, isAdmin, err := mysql.CheckUserCredentials(u.UserName, u.Password)
	if err != nil {
		if errors.Is(err, resp_code.ErrInvalidCredentials) {
			response.Code = resp_code.ErrorPwd
			zap.L().Warn("services-Login-CheckUserCredentials invalid credentials",
				zap.String("username", u.UserName))
			return
		} else {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-Login-CheckUserCredentials ", zap.Error(err))
			return
		}
	}

	// 生成token
	token, err := jwt.GenerateToken(userID, u.UserName, isAdmin)
	if err != nil {
		response.Code = resp_code.GenerateTokenError
		zap.L().Error("service-Login-GenerateToken ", zap.Error(err))
		return
	}
	// 设置响应的状态码和 token
	response.Code = resp_code.Success
	response.Data = token
	return
}

// GetUserDetail 获取用户详细信息
func (u *UserService) GetUserDetail() (response resp.ResponseWithData) {
	// 检验是否有这个用户ID
	exist, err := mysql.CheckUserID(u.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteUser-CheckUserID ", zap.Error(err))
		return
	}
	if !exist {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-DeleteUser-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", u.UserID), zap.Error(err))
		return
	}

	// 去数据库获取详细信息
	data, err := mysql.GetUserDetail(u.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-GetUserDetail-GetUserDetail "+
			fmt.Sprintf("db scan userID %d detail information failed ", u.UserID), zap.Error(err))
		return
	}
	response.Data = data
	response.Code = resp_code.Success
	return
}

// UpdateUserDetail 更新用户详细信息
func (u *UserService) UpdateUserDetail() (response resp.Response) {
	// 检验是否有这个用户ID
	exists, err := mysql.CheckUserID(u.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateUserDetail-CheckUserID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-UpdateUserDetail-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", u.UserID), zap.Error(err))
		return
	}
	if u.UserName != "" {
		// 检验是否有这个用户名
		exists, err = mysql.CheckUsername(u.UserName)
		if err != nil {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-UpdateUserDetail-CheckUsername ", zap.Error(err))
			return
		}
		if exists {
			response.Code = resp_code.UsernameAlreadyExist
			zap.L().Error("services-UpdateUserDetail-CheckUsername "+
				fmt.Sprintf("already have this username %s ", u.UserName), zap.Error(err))
			return
		}
	}

	// 更新用户信息
	if u.Email != "" {
		// 检查该邮箱是否已经存在
		exists, err = mysql.CheckEmail(u.Email)
		if err != nil {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-UpdateUserDetail-CheckUsername ", zap.Error(err))
			return
		}
		if exists {
			response.Code = resp_code.EmailAlreadyExist
			zap.L().Error("services-UpdateUserDetail-CheckUsername "+
				fmt.Sprintf("already have this username %s ", u.UserName), zap.Error(err))
			return
		}

		//验证码获取及验证
		code, err := redis.GetVerificationCode(u.Email)
		// 验证码过期
		if errors.Is(err, redis.Nil) {
			response.Code = resp_code.NeedObtainVerificationCode
			zap.L().Error("services-UpdateUserDetail-GetVerificationCode "+
				fmt.Sprintf("do not send verification code %s", u.Email), zap.Error(err))
			return
		}
		if errors.Is(err, fmt.Errorf("verification code expired")) {
			response.Code = resp_code.ExpiredVerCode
			zap.L().Error("services-UpdateUserDetail-GetVerificationCode "+
				fmt.Sprintf("expired verfiction code %s:%s ", u.Email, code), zap.Error(err))
			return
		}
		// 验证码错误
		if code != u.Code {
			response.Code = resp_code.ErrorVerCode
			zap.L().Error("services-UpdateUserDetail-GetVerificationCode "+
				fmt.Sprintf("error verfiction code %s:%s ", u.Email, code), zap.Error(err))
			return
		}
	}

	if u.Password != "" {
		// 密码加密
		u.Password, err = utils.CryptoPwd(u.Password)
		if err != nil {
			response.Code = resp_code.EncryptPwdError
			zap.L().Error("services-UpdateUserDetail-CryptoPwd ", zap.Error(err))
			return
		}
	}
	err = mysql.UpdateUserDetail(u.UserID, u.UserName, u.Email, u.Password)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateUserDetail-UpdateUserDetail "+
			fmt.Sprintf("db update userID %d failed ", u.UserID), zap.Error(err))
		return
	}
	response.Code = resp_code.Success
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
