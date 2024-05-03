package resp

import "online-judge/dao/mysql"

type Response struct {
	Code int `json:"code"`
}

// RegisterResponse 包含注册操作的返回结果
type RegisterResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

// SendCodeResponse 包含获取验证码的返回结果
type SendCodeResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

// GetDetailResponse 获取用户详细信息的返回结果
type GetDetailResponse struct {
	Code int        `json:"code"`
	Data mysql.User `json:"data"`
}

// DeleteUserResponse 删除用户的返回结果
type DeleteUserResponse struct {
	Code int `json:"code"`
}

// UpdateUserDetailResponse 更新用户详细信息的返回结果
type UpdateUserDetailResponse struct {
	Code int `json:"code"`
}
