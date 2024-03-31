package models

// RegisterResponse 包含注册操作的返回结果
type RegisterResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

// SendCodeResponse 包含注册操作的返回结果
type SendCodeResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}
