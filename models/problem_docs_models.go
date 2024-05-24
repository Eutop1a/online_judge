package models

type GetProblemListResponse struct {
	Code int `json:"code"` // "1000 获取题目列表成功" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type GetProblemDetailResponse struct {
	Code int `json:"code"` // "1000 获取题目列表成功" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type GetProblemIDResponse struct {
	Code int `json:"code"` // "1000 获取题目ID成功" "1020 题目title不存在" "1008 需要登录"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}

type GetProblemRandomResponse struct {
	Code int `json:"code"` // "1000 获取题目ID成功" "1008 需要登录"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}
