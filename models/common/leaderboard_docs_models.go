package common

type GetUserLeaderboardResponse struct {
	Code int `json:"code"` // "1000 获取用户题解排名成功" "1022 获取用户题解排名失败" "1014 服务器内部错误"

	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"` // omitempty 字段为空就忽略
}
