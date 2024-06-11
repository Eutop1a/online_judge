package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/models/common/response"
	"online_judge/models/user/request"
)

type ApiUser struct{}

// GetUserDetail 获取用户详细信息接口
// @Tags User API
// @Summary 获取用户详细信息
// @Description 获取用户详细信息接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} common.GetUserDetailResponse "获取用户信息成功"
// @Failure 200 {object} common.GetUserDetailResponse "参数错误"
// @Failure 200 {object} common.GetUserDetailResponse "没有此用户ID"
// @Failure 200 {object} common.GetUserDetailResponse "服务器内部错误"
// @Router /users/detail [GET]
func (u *ApiUser) GetUserDetail(c *gin.Context) {
	var req request.GetUserDetailReq
	uid, ok := c.Get(response.CtxUserIDKey)
	if !ok {
		response.ResponseError(c, response.CodeNeedLogin)
		return
	}
	//uid := c.Param("user_id")
	//if uid == "" {
	//	zap.L().Error("GetUserDetail params error")
	//	response.ResponseError(c, response.CodeInvalidParam)
	//	return
	//}
	//getDetail.UserID, _ = strconv.ParseInt(uid, 10, 64)
	req.UserID = uid.(int64)
	var ret response.ResponseWithData
	ret = UserService.GetUserDetail(req)

	switch ret.Code {
	// 成功
	case resp_code.Success:
		response.ResponseSuccess(c, ret.Data)

	// 用户不存在
	case resp_code.NotExistUserID:
		response.ResponseError(c, response.CodeUseNotExist)

	// 内部错误
	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// UpdateUserDetail 更新用户详细信息接口
// @Tags User API
// @Summary 更新用户详细信息
// @Description 更新用户详细信息接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param username formData string false "用户名"
// @Param password formData string false "用户密码"
// @Param email formData string false "用户邮箱"
// @Param code formData string false "邮箱验证码"
// @Success 200 {object} common.UpdateUserDetailResponse "更新用户信息成功"
// @Failure 200 {object} common.UpdateUserDetailResponse "参数错误"
// @Failure 200 {object} common.UpdateUserDetailResponse "没有此用户ID"
// @Failure 200 {object} common.UpdateUserDetailResponse "验证码错误"
// @Failure 200 {object} common.UpdateUserDetailResponse "验证码过期"
// @Failure 200 {object} common.UpdateUserDetailResponse "服务器内部错误"
// @Router /users/update [PUT]
func (u *ApiUser) UpdateUserDetail(c *gin.Context) {
	var req request.UpdateUserDetailReq
	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("UpdateUserDetail.ShouldBind error " + err.Error())
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	uid, ok := c.Get(response.CtxUserIDKey)
	if !ok {
		response.ResponseError(c, response.CodeNeedLogin)
		return
	}
	//uid := c.Param("user_id")
	//if uid == "" {
	//	zap.L().Error("UpdateUserDetail params error")
	//	response.ResponseError(c, response.CodeInvalidParam)
	//	return
	//}
	//update.UserID, _ = strconv.ParseInt(uid, 10, 64)
	//fmt.Println(update.UserID)
	req.UserID = uid.(int64)
	var ret response.Response
	ret = UserService.UpdateUserDetail(req)

	switch ret.Code {
	// 成功
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	// 用户不存在
	case resp_code.NotExistUserID:
		response.ResponseError(c, response.CodeUseNotExist)

	// 验证码错误
	case resp_code.ErrorVerCode:
		response.ResponseError(c, response.CodeErrorVerCode)

	// 验证码过期
	case resp_code.ExpiredVerCode:
		response.ResponseError(c, response.CodeExpiredVerCode)

	// 新用户名已经存在
	case resp_code.UsernameAlreadyExist:
		response.ResponseError(c, response.CodeUsernameAlreadyExist)

	// 未申请验证码
	case resp_code.NeedObtainVerificationCode:
		response.ResponseError(c, response.CodeObtainVerificationCode)

	// 邮箱已经存在
	case resp_code.EmailAlreadyExist:
		response.ResponseError(c, response.CodeEmailExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// GetUserID 获取用户ID接口
// @Tags User API
// @Summary 获取用户ID
// @Description 获取用户ID接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param username formData string true "用户名"
// @Success 200 {object} common.GetUserIDResponse "获取用户ID成功"
// @Failure 200 {object} common.GetUserIDResponse "用户名不存在"
// @Router /users/user-id [POST]
func (u *ApiUser) GetUserID(c *gin.Context) {
	var req request.GetUserIDReq
	req.Username = c.PostForm("username")
	uid, err := UserService.GetUserID(req)
	if err != nil {
		response.ResponseError(c, response.CodeUseNotExist)
		return
	}
	response.ResponseSuccess(c, uid)
}
