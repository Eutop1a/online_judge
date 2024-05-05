package services

import (
	"fmt"
	"go.uber.org/zap"
	"online-judge/dao/mysql"
	"online-judge/pkg/resp"
	"time"
)

// Submission 提交记录
type Submission struct {
	SubmissionID   string    `form:"submission_id" json:"submission_id"`     // 提交ID
	UserID         int64     `form:"user_id" json:"user_id"`                 //用户ID
	ProblemID      string    `form:"problem_id" json:"problem_id"`           //题目ID
	Language       string    `form:"language" json:"language"`               //编程语言
	Code           string    `form:"code" json:"code"`                       // 代码
	SubmissionTime time.Time `form:"submission_time" json:"submission_time"` // 提交时间
}

func (s *Submission) SubmitCode() (response resp.Response) {
	// 检查用户id是否存在
	var idNum int64
	err := mysql.CheckUserID(s.UserID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-CheckUserID ", zap.Error(err))
		return
	case idNum == 0: // 用户id不存在
		response.Code = resp.NotExistUserID
		zap.L().Error("services-SubmitCode-CheckUserID " +
			fmt.Sprintf(" user_id %d do not exist ", s.UserID))
		return
	}
	idNum = 0
	// 检查题目id是否存在
	err = mysql.CheckProblemID(s.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-CheckProblemID ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = resp.ProblemNotExist
		zap.L().Error("services-SubmitCode-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist ", s.ProblemID))
		return
	}
	tmp := mysql.Submission{
		UserID:         s.UserID,
		SubmissionID:   s.SubmissionID,
		ProblemID:      s.ProblemID,
		Language:       s.Language,
		Code:           s.Code,
		SubmissionTime: s.SubmissionTime,
	}

	err = mysql.SubmitCode(&tmp)
	if err != nil {
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-SubmitCode  ", zap.Error(err))
		return
	}
	response.Code = resp.Success
	return
}
