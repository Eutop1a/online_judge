package models

type _Response struct {
	Code int `json:"code"` // 状态码，不同的状态码对应不同的结果：

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type AddSuperAdminResponse struct {
	Code int `json:"code"` // "1000 添加超级管理员成功" "1001 参数错误" "1005 没有此用户ID" "1026 用户已是管理员" "1025 密钥错误" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type DeleteUserResponse struct {
	Code int `json:"code"` // "1000 删除用户成功" "1001 参数错误" "1004 没有此用户ID" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type AddAdminResponse struct {
	Code int `json:"code"` // "1000 删除用户成功" "1001 参数错误" "1005 没有此用户ID" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type CreateProblemResponse struct {
	Code int `json:"code"` // "1000 创建成功" "1001 参数错误" "1018 测试用例格式错误" "1019 题目标题已存在" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type UpdateProblemResponse struct {
	Code int `json:"code"` // "1000 修改成功" "1021 题目ID不存在" "1019 题目标题已存在" "1018 测试用例格式错误" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

type DeleteProblemResponse struct {
	Code int `json:"code"` // "1000 删除成功" "1021 题目ID不存在" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
