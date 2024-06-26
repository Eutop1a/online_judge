package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"strconv"
)

type ApiAdminUser struct{}

// AddSuperAdmin 添加超级管理员接口
// @Tags Admin API
// @Summary 添加超级管理员
// @Description 添加超级管理员接口
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param secret formData string true "密钥"
// @Success 200 {object} common.AddSuperAdminResponse "添加超级管理员成功"
// @Failure 200 {object} common.AddSuperAdminResponse "参数错误"
// @Failure 200 {object} common.AddSuperAdminResponse "没有此用户ID"
// @Failure 200 {object} common.AddSuperAdminResponse "用户已是管理员"
// @Failure 200 {object} common.AddSuperAdminResponse "密钥错误"
// @Failure 200 {object} common.AddSuperAdminResponse "服务器内部错误"
// @Router /admin/users/add-super-admin [POST]
func (a *ApiAdminUser) AddSuperAdmin(c *gin.Context) {
	var addSuperAdminReq request.AdminAddSuperAdminReq
	//uid := c.PostForm("user_id")
	if err := c.ShouldBind(&addSuperAdminReq); err != nil {
		zap.L().Error("controller-ShouldBind ", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	if addSuperAdminReq.Username == "" {
		zap.L().Error("controller-AddSuperAdmin-PostForm add admin params error")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	var ret response.Response
	ret = AdminService.AddSuperAdmin(addSuperAdminReq)
	switch ret.Code {

	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.UserAlreadyRoot:
		response.ResponseError(c, response.CodeUserAlreadyRoot)

	case resp_code.UsernameDoesNotExist:
		response.ResponseError(c, response.CodeUsernameNotExist)

	case resp_code.SecretError:
		response.ResponseError(c, response.CodeErrorSecret)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
	return
}

// DeleteUser 删除用户接口
// @Tags Admin API
// @Summary 删除用户
// @Description 删除用户接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param user_id path string true "用户ID"
// @Success 200 {object} common.DeleteUserResponse "删除用户成功"
// @Failure 200 {object} common.DeleteUserResponse "参数错误"
// @Failure 200 {object} common.DeleteUserResponse "没有此用户ID"
// @Failure 200 {object} common.DeleteUserResponse "需要登录"
// @Failure 200 {object} common.DeleteUserResponse "服务器内部错误"
// @Router /admin/users/{user_id} [DELETE]
func (a *ApiAdminUser) DeleteUser(c *gin.Context) {
	var deleteUserReq request.AdminDeleteUserReq
	uid := c.Param("user_id")
	if uid == "" {
		zap.L().Error("controller-DeleteUser-Param deleteUser params error")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	deleteUserReq.UserID, _ = strconv.ParseInt(uid, 10, 64)

	ret := AdminService.DeleteUser(deleteUserReq)

	switch ret.Code {
	// 成功
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	// 用户不存在
	case resp_code.NotExistUserID:
		response.ResponseError(c, response.CodeUseNotExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}
