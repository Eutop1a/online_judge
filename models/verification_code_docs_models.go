package models

type SendEmailCodeResponse struct {
	Code int `json:"code"` //"1000 发送邮箱验证码成功" "1015 邮箱格式错误" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type SendPictureCodeResponse struct {
	Code int `json:"code"` //"1000 发送图片验证码成功" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}
type CheckPictureCodeResponse struct {
	Code int `json:"code"` //"1000 图片验证码正确" "1017 图片验证码错误" "1014 用户名不存在"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}
