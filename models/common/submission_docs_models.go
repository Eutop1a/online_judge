package common

type SubmitCodeResponse struct {
	Code int `json:"code"` // "1000 提交代码成功" "1005 用户ID不存在" "1021 题目ID不存在" "1024 不支持的语言类型" "1008 需要登录" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}
