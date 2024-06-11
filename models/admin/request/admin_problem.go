package request

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// AdminCreateProblemReq 创建题目
type AdminCreateProblemReq struct {
	MaxRuntime int         `form:"max_runtime" json:"max_runtime"` // 时间限制
	MaxMemory  int         `form:"max_memory" json:"max_memory"`   // 内存限制
	ProblemID  string      `form:"problem_id" json:"problem_id"`   // unique key
	Title      string      `form:"title" json:"title"`             // problem title
	Content    string      `form:"content" json:"content"`         // problem description
	Difficulty string      `form:"difficulty" json:"difficulty"`   // easy mid hard
	TestCases  []*TestCase `form:"test_cases" json:"test_cases"`   // 测试样例集

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}

// TestCase 测试样例
type TestCase struct {
	TID      string `form:"TID" json:"TID"`           // testCase ID
	PID      string `form:"PID" json:"PID"`           // 对应的题目ID
	Input    string `form:"input" json:"input"`       // 输入
	Expected string `form:"expected" json:"expected"` // 期望输出
}

// AdminUpdateProblemReq 更新题目
type AdminUpdateProblemReq struct {
	MaxRuntime int         `form:"max_runtime" json:"max_runtime"` // 时间限制
	MaxMemory  int         `form:"max_memory" json:"max_memory"`   // 内存限制
	ProblemID  string      `form:"problem_id" json:"problem_id"`   // unique key
	Title      string      `form:"title" json:"title"`             // problem title
	Content    string      `form:"content" json:"content"`         // problem description
	Difficulty string      `form:"difficulty" json:"difficulty"`   // easy mid hard
	TestCases  []*TestCase `form:"test_cases" json:"test_cases"`   // 测试样例集

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}

// AdminDeleteProblemReq 删除题目
type AdminDeleteProblemReq struct {
	ProblemID string `form:"problem_id" json:"problem_id"` // unique key

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}

// AdminCreateProblemWithFileReq 创建测试样例为文件的题目
type AdminCreateProblemWithFileReq struct {
	MaxRuntime        int                 `form:"max_runtime" json:"max_runtime"`                   // 时间限制
	MaxMemory         int                 `form:"max_memory" json:"max_memory"`                     // 内存限制
	ProblemID         string              `form:"problem_id" json:"problem_id"`                     // unique key
	Title             string              `form:"title" json:"title"`                               // problem title
	Content           string              `form:"content" json:"content"`                           // problem description
	Difficulty        string              `form:"difficulty" json:"difficulty"`                     // easy mid hard
	InputDst          string              `form:"input_dst" json:"input_dst"`                       // 输入文件保存的地址
	ExpectedDst       string              `form:"expected_dst" json:"expected_dst"`                 // 输出文件保存的地址
	TestCasesWithFile []*TestCaseWithFile `form:"test_cases_with_file" json:"test_cases_with_file"` // 测试样例集(file)

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}

// TestCaseWithFile 文件测试样例
type TestCaseWithFile struct {
	TID          string `form:"tid" json:"tid"`                     // 测试样例ID
	PID          string `form:"pid" json:"pid"`                     // 对应的题目ID
	InputPath    string `form:"input_path" json:"input_path"`       // 输入文件
	ExpectedPath string `form:"expected_path" json:"expected_path"` // 期望输出文件名
}

// AdminDeleteProblemWithFileReq 删除测试样例为文件的题目
type AdminDeleteProblemWithFileReq struct {
	ProblemID string `form:"problem_id" json:"problem_id"` // unique key
	Title     string `form:"title" json:"title"`           // problem title

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}

// AdminDeleteProblemTestCaseWithFileReq 删除测试样例为文件的测试样例
type AdminDeleteProblemTestCaseWithFileReq struct {
	ProblemID string `form:"problem_id" json:"problem_id"` // unique key
	Title     string `form:"title" json:"title"`           // problem title

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}

// AdminUpdateProblemWithFileReq 更新测试样例为文件的题目
type AdminUpdateProblemWithFileReq struct {
	MaxRuntime        int                 `form:"max_runtime" json:"max_runtime"`                   // 时间限制
	MaxMemory         int                 `form:"max_memory" json:"max_memory"`                     // 内存限制
	ProblemID         string              `form:"problem_id" json:"problem_id"`                     // unique key
	Title             string              `form:"title" json:"title"`                               // problem title
	Content           string              `form:"content" json:"content"`                           // problem description
	Difficulty        string              `form:"difficulty" json:"difficulty"`                     // easy mid hard
	InputDst          string              `form:"input_dst" json:"input_dst"`                       // 输入文件保存的地址
	ExpectedDst       string              `form:"expected_dst" json:"expected_dst"`                 // 输出文件保存的地址
	TestCasesWithFile []*TestCaseWithFile `form:"test_cases_with_file" json:"test_cases_with_file"` // 测试样例集(file)

	RedisClient *redis.Client   `form:"redis_client" json:"redis_client"`
	Ctx         context.Context `form:"context" json:"context"`
}
