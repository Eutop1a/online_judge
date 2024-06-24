package user

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/dao/redis/cache/verify"
	"online_judge/models/common/response"
	"online_judge/models/user/request"
	"online_judge/pkg/utils"
)

type UserService struct{}

// GetUserDetail 获取用户详细信息
func (u *UserService) GetUserDetail(req request.GetUserDetailReq) (response response.ResponseWithData) {
	// 检验是否有这个用户ID
	exist, err := mysql.CheckUserID(req.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteUser-CheckUserID ", zap.Error(err))
		return
	}
	if !exist {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-DeleteUser-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", req.UserID), zap.Error(err))
		return
	}

	// 去数据库获取详细信息
	data, err := mysql.GetUserDetail(req.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-GetUserDetail-GetUserDetail "+
			fmt.Sprintf("db scan userID %d detail information failed ", req.UserID), zap.Error(err))
		return
	}
	response.Data = data
	response.Code = resp_code.Success
	return
}

// UpdateUserDetail 更新用户详细信息
func (u *UserService) UpdateUserDetail(req request.UpdateUserDetailReq) (response response.Response) {
	// 检验是否有这个用户ID
	exists, err := mysql.CheckUserID(req.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateUserDetail-CheckUserID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-UpdateUserDetail-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", req.UserID), zap.Error(err))
		return
	}
	if req.Username != "" {
		// 检验是否有这个用户名
		exists, err = mysql.CheckUsername(req.Username)
		if err != nil {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-UpdateUserDetail-CheckUsername ", zap.Error(err))
			return
		}
		if exists {
			response.Code = resp_code.UsernameAlreadyExist
			zap.L().Error("services-UpdateUserDetail-CheckUsername "+
				fmt.Sprintf("already have this username %s ", req.Username), zap.Error(err))
			return
		}
	}

	// 更新用户信息
	if req.Email != "" {
		// 检查该邮箱是否已经存在
		exists, err = mysql.CheckEmail(req.Email)
		if err != nil {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-UpdateUserDetail-CheckUsername ", zap.Error(err))
			return
		}
		if exists {
			response.Code = resp_code.EmailAlreadyExist
			zap.L().Error("services-UpdateUserDetail-CheckUsername "+
				fmt.Sprintf("already have this username %s ", req.Username), zap.Error(err))
			return
		}

		//验证码获取及验证
		code, err := verify.GetVerifyCode(req.Email)
		// 验证码过期
		if err == redis.Nil {
			response.Code = resp_code.NeedObtainVerificationCode
			zap.L().Error("services-UpdateUserDetail-GetVerifyCode "+
				fmt.Sprintf("do not send verify code %s", req.Email), zap.Error(err))
			return
		}
		if err == fmt.Errorf("verify code expired") {
			response.Code = resp_code.ExpiredVerCode
			zap.L().Error("services-UpdateUserDetail-GetVerifyCode "+
				fmt.Sprintf("expired verfiction code %s:%s ", req.Email, code), zap.Error(err))
			return
		}
		// 验证码错误
		if code != req.Code {
			response.Code = resp_code.ErrorVerCode
			zap.L().Error("services-UpdateUserDetail-GetVerifyCode "+
				fmt.Sprintf("error verfiction code %s:%s ", req.Email, code), zap.Error(err))
			return
		}
	}

	if req.Password != "" {
		// 密码加密
		req.Password, err = utils.CryptoPwd(req.Password)
		if err != nil {
			response.Code = resp_code.EncryptPwdError
			zap.L().Error("services-UpdateUserDetail-CryptoPwd ", zap.Error(err))
			return
		}
	}
	err = mysql.UpdateUserDetail(req.UserID, req.Username, req.Email, req.Password)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateUserDetail-UpdateUserDetail "+
			fmt.Sprintf("db update userID %d failed ", req.UserID), zap.Error(err))
		return
	}
	response.Code = resp_code.Success
	return
}

func (u *UserService) GetUserID(req request.GetUserIDReq) (uid int64, err error) {
	uid, err = mysql.GetUserID(req.Username)
	if err != nil {
		zap.L().Error("services-GetUserID-GetUserID "+req.Username+"err: ", zap.Error(err))
		return 0, err
	}
	return uid, nil
}
