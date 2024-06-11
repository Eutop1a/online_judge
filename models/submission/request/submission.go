package request

import "time"

// SubmissionReq 提交记录
type SubmissionReq struct {
	SubmissionID   string    `form:"submission_id" json:"submission_id"`     // 提交ID
	UserID         int64     `form:"user_id" json:"user_id"`                 //用户ID
	ProblemID      string    `form:"problem_id" json:"problem_id"`           //题目ID
	Language       string    `form:"language" json:"language"`               //编程语言
	Code           string    `form:"code" json:"code"`                       // 代码
	SubmissionTime time.Time `form:"submission_time" json:"submission_time"` // 提交时间
}
