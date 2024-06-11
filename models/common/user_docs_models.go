package common

type RegisterResponse struct {
	Code int `json:"code"` // "1000 注册成功" "1002 用户已存在" "1011 验证码错误或已过期" "1012 验证码过期" "1013 该邮箱已经存在" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type LoginResponse struct {
	Code int `json:"code"` // "1000 登录成功" "1001 参数错误" "1004 用户名不存在" "1004 验证码错误" "1011 验证码过期" "1006 密码错误" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type GetUserDetailResponse struct {
	Code int `json:"code"` //"1000 获取用户信息成功" "1001 参数错误" "1004 没有此用户ID" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type UpdateUserDetailResponse struct {
	Code int `json:"code"` //"1000 更新用户信息成功" "1001 参数错误" "1004 没有此用户ID" "1011 验证码错误" "1012 验证码过期" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type GetUserIDResponse struct {
	Code int `json:"code"` //"1000 获取用户ID成功" "1004 用户名不存在"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}
