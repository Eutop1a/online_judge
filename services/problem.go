package services

import (
	"go.uber.org/zap"
	"online-judge/dao/mysql"
)

// Problem 问题结构体
type Problem struct {
	ID                int                 `form:"id" json:"id"`                                     // primary key
	MaxRuntime        int                 `form:"max_runtime" json:"max_runtime"`                   // 时间限制
	MaxMemory         int                 `form:"max_memory" json:"max_memory"`                     // 内存限制
	Size              int                 `form:"size" json:"size"`                                 // 每页的记录数
	Page              int                 `form:"page" json:"page"`                                 // 第page页
	Count             int64               `form:"count" json:"count"`                               // 查到的记录数
	ProblemID         string              `form:"problem_id" json:"problem_id"`                     // unique key
	Title             string              `form:"title" json:"title"`                               // problem title
	Content           string              `form:"content" json:"content"`                           // problem description
	Difficulty        string              `form:"difficulty" json:"difficulty"`                     // easy mid hard
	InputDst          string              `form:"input_dst" json:"input_dst"`                       // 输入文件保存的地址
	ExpectedDst       string              `form:"expected_dst" json:"expected_dst"`                 // 输出文件保存的地址
	TestCases         []*TestCase         `form:"test_cases" json:"test_cases"`                     // 测试样例集
	TestCasesWithFile []*TestCaseWithFile `form:"test_cases_with_file" json:"test_cases_with_file"` // 测试样例集(file)

}
type TestCaseWithFile struct {
	TID          string `form:"tid" json:"tid"`                     // 测试样例ID
	PID          string `form:"pid" json:"pid"`                     // 对应的题目ID
	InputPath    string `form:"input_path" json:"input_path"`       // 输入文件
	ExpectedPath string `form:"expected_path" json:"expected_path"` // 期望输出文件名
}

// TestCase 测试样例
type TestCase struct {
	TID      string `form:"TID" json:"TID"`           // testCase ID
	PID      string `form:"PID" json:"PID"`           // 对应的题目ID
	Input    string `form:"input" json:"input"`       // 输入
	Expected string `form:"expected" json:"expected"` // 期望输出
}

// GetProblemList 获取题目列表
func (p *Problem) GetProblemList() (*[]Problem, error) {
	var count int64
	data, err := mysql.GetProblemList(p.Page, p.Size, &count)
	if err != nil {
		zap.L().Error("services-GetProblemList-GetProblemList ", zap.Error(err))
		return nil, err
	}

	problems := make([]Problem, len(*data))
	for k, v := range *data {
		problems[k].ID = v.ID
		problems[k].ProblemID = v.ProblemID
		problems[k].Content = v.Content
		problems[k].Title = v.Title
		problems[k].Difficulty = v.Difficulty
		problems[k].MaxMemory = v.MaxMemory
		problems[k].MaxRuntime = v.MaxRuntime
		problems[k].MaxRuntime = v.MaxRuntime
		problems[k].Count = count
		problems[k].Size = p.Size
		problems[k].Page = p.Page
		problems[k].TestCases = make([]*TestCase, len(v.TestCases))
		for i, tc := range v.TestCases {
			problems[k].TestCases[i] = &TestCase{
				TID:      tc.TID,
				PID:      tc.PID,
				Input:    tc.Input,
				Expected: tc.Expected,
			}
		}
	}

	return &problems, nil
}

// GetProblemDetail 获取单个题目详细信息
func (p *Problem) GetProblemDetail() (*mysql.Problems, error) {
	data, err := mysql.GetProblemDetail(p.ProblemID)
	if err != nil {
		zap.L().Error("services-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return nil, err
	}
	return data, nil
}

// GetProblemID 获取题目ID
func (p *Problem) GetProblemID() (problemID string, err error) {
	problemID, err = mysql.GetProblemID(p.Title)
	if err != nil {
		zap.L().Error("services-GetProblemID-GetProblemID", zap.Error(err))
		return "", err
	}
	return
}

// GetProblemRandom 随机获取一个题目
func (p *Problem) GetProblemRandom() (*mysql.Problems, error) {
	problem, err := mysql.GetProblemRandom()
	if err != nil {
		zap.L().Error("services-GetProblemRandom-GetProblemRandom", zap.Error(err))
		return nil, err
	}
	return problem, nil
}
