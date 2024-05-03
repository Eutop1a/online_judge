package services

import (
	"fmt"
	"go.uber.org/zap"
	"online-judge/dao/mysql"
	"online-judge/pkg/resp"
)

type Problem struct {
	ProblemID  string      `form:"problem_id" json:"problem_id"`   // primary key
	Title      string      `form:"title" json:"title"`             // problem title
	Content    string      `form:"content" json:"content"`         // problem description
	Difficulty string      `form:"difficulty" json:"difficulty"`   // easy mid hard
	MaxRuntime int         `form:"max_runtime" json:"max_runtime"` // 时间限制
	MaxMemory  int         `form:"max_memory" json:"max_memory"`   // 内存限制
	TestCases  []*TestCase `form:"test_cases" json:"test_cases"`   // 测试样例集
}

// TestCase 测试样例
type TestCase struct {
	TID      string `form:"TID" json:"TID"`           // testCase ID
	PID      string `form:"PID" json:"PID"`           // 对应的题目ID
	Input    string `form:"input" json:"input"`       // 输入
	Expected string `form:"expected" json:"expected"` // 期望输出
}

func (p *Problem) GetProblemList() (*[]mysql.Problems, error) {
	data, err := mysql.GetProblemList()
	if err != nil {
		zap.L().Error("mysql.GetProblemList() failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func (p *Problem) GetProblemDetail() (*mysql.Problems, error) {
	data, err := mysql.GetProblemDetail(p.ProblemID)
	if err != nil {
		zap.L().Error("mysql.GetProblemDetail() error", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func (p *Problem) CreateProblem() (response resp.RegisterResponse) {
	// 检查题目标题是否重复
	var problemNum int64
	err := mysql.CheckProblemTitle(p.Title, &problemNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp.SearchDBError
		zap.L().Error("services-SearchDBError", zap.Error(err))
		return
	case problemNum > 0: // 题目已经存在
		response.Code = resp.ProblemAlreadyExist
		zap.L().Error("services-" + fmt.Sprintf("Title %s aleardy exist", p.Title))
		return
	}
	var problems mysql.Problems
	problems.ProblemID = p.ProblemID
	problems.Title = p.Title
	problems.Content = p.Content
	problems.Difficulty = p.Difficulty
	problems.MaxRuntime = p.MaxRuntime
	problems.MaxMemory = p.MaxMemory

	// 提前转换类型
	var convertedTestCases []*mysql.TestCase
	for _, tc := range p.TestCases {
		// 进行类型转换
		convertedTestCases = append(convertedTestCases, &mysql.TestCase{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		})
	}

	problems.TestCases = convertedTestCases
	// 创建题目
	err = mysql.CreateProblem(&problems)
	if err != nil {
		response.Code = resp.CreateProblemError
		zap.L().Error("services-SearchDBError", zap.Error(err))
		return
	}
	response.Code = resp.Success
	return
}

func (p *Problem) UpdateProblem() {

}

func (p *Problem) DeleteProblem() {

}
