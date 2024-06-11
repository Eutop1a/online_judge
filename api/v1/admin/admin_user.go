package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"strconv"
)

type ApiAdmin struct{}

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
func (a *ApiAdmin) AddSuperAdmin(c *gin.Context) {
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
		response.ResponseSuccess(c, resp_code.Success)

	case resp_code.NotExistUsername:
		response.ResponseError(c, response.CodeUsernameNotExist)

	case resp_code.SecretError:
		response.ResponseError(c, response.CodeErrorSecret)

	case resp_code.UsernameAlreadyExist:
		response.ResponseError(c, response.CodeUsernameAlreadyExist)

	default:
		response.ResponseError(c, response.CodeInvalidParam)
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
func (a *ApiAdmin) DeleteUser(c *gin.Context) {
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

// AddAdmin 添加管理员接口
// @Tags Admin API
// @Summary 添加管理员
// @Description 添加管理员接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param username formData string true "用户名"
// @Success 200 {object} common.AddAdminResponse "删除用户成功"
// @Failure 200 {object} common.AddAdminResponse "参数错误"
// @Failure 200 {object} common.AddAdminResponse "没有此用户ID"
// @Failure 200 {object} common.AddAdminResponse "需要登录"
// @Failure 200 {object} common.AddAdminResponse "服务器内部错误"
// @Router /admin/users/add-admin [POST]
func (a *ApiAdmin) AddAdmin(c *gin.Context) {
	var req request.AdminAddAdminReq
	//uid := c.PostForm("user_id")
	req.Username = c.PostForm("username")
	if req.Username == "" {
		zap.L().Error("add admin params error")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	//addAdmin.UserID, _ = strconv.ParseInt(uid, 10, 64)

	var ret response.Response
	ret = AdminService.AddAdmin(req)

	switch ret.Code {
	// 成功
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	// 用户名不存在
	case resp_code.NotExistUsername:
		response.ResponseError(c, response.CodeUsernameNotExist)

	// 用户已经是管理员了
	case resp_code.UsernameAlreadyExist:
		response.ResponseError(c, response.CodeUsernameAlreadyExist)

	// 服务器内部错误
	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}
