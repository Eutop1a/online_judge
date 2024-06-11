package admin

import (
	"fmt"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"online_judge/pkg/utils"
)

type AdminUserService struct{}

// DeleteUser 删除用户
func (u *AdminUserService) DeleteUser(request request.AdminDeleteUserReq) (response response.Response) {
	// 检验是否有这个用户ID
	exist, err := mysql.CheckUserID(request.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteUser-CheckUserID ", zap.Error(err))
		return
	}
	if !exist {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-DeleteUser-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", request.UserID), zap.Error(err))
		return
	}
	// 删除用户
	err = mysql.DeleteUser(request.UserID)
	if err != nil {
		response.Code = resp_code.DBDeleteError
		zap.L().Error("services-DeleteUser-DeleteUser "+
			fmt.Sprintf("delete userID %d failed ", request.UserID), zap.Error(err))
		return
	}
	// 删除成功
	response.Code = resp_code.Success
	return
}

const SECRETCIPHER = "afd372788c1f7f646a678654901ce041ecc9012487dc0055b932cac9acaca27b6cf0488a3b5d0aa05022ab9a51e54b0e54e8188beaf4ad9cef517c0c76641f21"

func (u *AdminUserService) AddSuperAdmin(request request.AdminAddSuperAdminReq) (response response.Response) {
	if utils.CryptoSecret(request.Secret) != SECRETCIPHER {
		response.Code = resp_code.SecretError
		zap.L().Error("services-AddSuperAdmin-CryptoSecret " +
			fmt.Sprintf("secret error %s:%s", request.Username, request.Secret))
		return
	}

	// 检查改用户名是否已经存在, 已经存在后是否为管理员
	userExists, adminExists, err := mysql.CheckUsernameAndAdminExists(request.Username)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-AddSuperAdmin-CheckUsernameAndAdminExists ", zap.Error(err))
		return
	}
	if !userExists {
		response.Code = resp_code.NotExistUsername
		zap.L().Error("services-AddSuperAdmin-CheckUsername "+
			fmt.Sprintf("do not have this username %d ", request.Username), zap.Error(err))
		return
	}
	if adminExists {
		response.Code = resp_code.UsernameAlreadyExist
		zap.L().Error("services-AddSuperAdmin-CheckUsernameAlreadyExists "+
			fmt.Sprintf("already have this username %s ", request.Username), zap.Error(err))
		return
	}
	// 如果不是管理员就添加到数据库中
	err = mysql.AddAdminUserByUsername(request.Username)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-AddAdmin-AddAdminUser ", zap.Error(err))
		return
	}
	response.Code = resp_code.Success

	return
}

// AddAdmin 添加管理员
func (u *AdminUserService) AddAdmin(request request.AdminAddAdminReq) (response response.Response) {
	// 检查改用户名是否已经存在已经存在后是否为管理员
	userExists, adminExists, err := mysql.CheckUsernameAndAdminExists(request.Username)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-AddSuperAdmin-CheckUsernameAndAdminExists ", zap.Error(err))
		return
	}
	if !userExists {
		response.Code = resp_code.NotExistUsername
		zap.L().Error("services-AddSuperAdmin-CheckUsername "+
			fmt.Sprintf("do not have this username %d ", request.Username), zap.Error(err))
		return
	}
	if adminExists {
		response.Code = resp_code.UsernameAlreadyExist
		zap.L().Error("services-AddSuperAdmin-CheckUsernameAlreadyExists "+
			fmt.Sprintf("already have this username %s ", request.Username), zap.Error(err))
		return
	}

	err = mysql.AddAdminUserByUsername(request.Username)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-AddAdmin-AddAdminUser ", zap.Error(err))
		return
	}
	response.Code = resp_code.Success

	return
}
