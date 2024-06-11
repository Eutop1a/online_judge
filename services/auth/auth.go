package auth

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/dao/redis"
	"online_judge/models/auth/request"
	"online_judge/models/common/response"
	"online_judge/pkg/jwt"
	"online_judge/pkg/snowflake"
	"online_judge/pkg/utils"
)

type AuthService struct{}

// Register 用户注册服务
func (a *AuthService) Register(request request.AuthRegisterReq) (response response.Response) {
	// 检查数据库中是否已经有这个email或username
	var countEmail, countUsername int64
	err := mysql.CheckEmailAndUsername(request.Email, request.UserName, &countEmail, &countUsername)
	// UID作为唯一标识，但是email和username不能重复
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-Register-CheckEmailAndUsername ", zap.Error(err))
		return
	}

	if countEmail > 0 {
		response.Code = resp_code.EmailAlreadyExist
		zap.L().Warn("services-Register-CheckEmailAndUsername email already exist",
			zap.String("email", request.Email))
		return
	}

	if countUsername > 0 {
		response.Code = resp_code.UsernameAlreadyExist
		zap.L().Warn("services-Register-CheckEmailAndUsername username already exist",
			zap.String("username", request.UserName))
		return
	}

	//验证码获取及验证
	code, err := redis.GetVerificationCode(request.Email)
	// 验证码过期
	if err != nil {
		if errors.Is(err, fmt.Errorf("verify code expired")) {
			response.Code = resp_code.ExpiredVerCode
			zap.L().Warn("services-Register-GetVerificationCode verify code expired",
				zap.String("email", request.Email))
		} else {
			response.Code = resp_code.ErrorVerCode
			zap.L().Warn("services-Register-GetVerificationCode error verify code",
				zap.String("email", request.Email))
		}
		return
	}
	// 验证码错误
	if code != request.Code {
		response.Code = resp_code.ErrorVerCode
		zap.L().Warn("services-Register-GetVerificationCode incorrect verify code",
			zap.String("email", request.Email))
		return
	}

	// 生成唯一的ID
	Id := snowflake.GetID()
	request.UserID = Id
	// 密码加密
	request.Password, err = utils.CryptoPwd(request.Password)
	if err != nil {
		response.Code = resp_code.EncryptPwdError
		zap.L().Error("services-Register-CryptoPwd ", zap.Error(err))
		return
	}

	// 插入数据库
	if err = mysql.InsertNewUser(request.UserID, request.UserName, request.Password, request.Email); err != nil {
		response.Code = resp_code.InsertNewUserError
		zap.L().Error("services-Register-InsertNewUser "+
			fmt.Sprintf("insert new user %s error ", request.UserName), zap.Error(err))
		return
	}
	// 设置响应的状态码
	response.Code = resp_code.Success
	return
}

// Login 用户登录服务
func (a *AuthService) Login(request request.AuthLoginReq) (response response.ResponseWithData) {
	// 检查用户名和密码是否正确
	userID, isAdmin, err := mysql.CheckUserCredentials(request.UserName, request.Password)
	if err != nil {
		if errors.Is(err, resp_code.ErrInvalidCredentials) {
			response.Code = resp_code.ErrorPwd
			zap.L().Warn("services-Login-CheckUserCredentials invalid credentials",
				zap.String("username", request.UserName))
			return
		} else {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-Login-CheckUserCredentials ", zap.Error(err))
			return
		}
	}

	// 生成token
	token, err := jwt.GenerateToken(userID, request.UserName, isAdmin)
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
