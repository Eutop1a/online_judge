package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"online-judge/dao/mysql"
	"online-judge/pkg/define"
	"time"
)

type TestCaseTmp struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type CreateProblemRequest struct {
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	Difficulty string        `json:"difficulty"`
	MaxRuntime int           `json:"max_runtime"`
	MaxMemory  int           `json:"max_memory"`
	TestCases  []TestCaseTmp `json:"test_cases"`
}

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

var (
	Nil = redis.Nil
	Ctx = context.Background()
)

// GetProblemList 获取题目列表
func (p *Problem) GetProblemList() (*[]Problem, error) {
	var count int64
	problemList, err := mysql.GetProblemList(p.Page, p.Size, &count)
	if err != nil {
		zap.L().Error("services-GetProblemList-GetProblemList ", zap.Error(err))
		return nil, err
	}

	problems := make([]Problem, len(problemList))
	for k, v := range problemList {
		problems[k].ID = v.ID
		problems[k].ProblemID = v.ProblemID
		problems[k].Title = v.Title
		problems[k].Difficulty = v.Difficulty
		problems[k].Count = count
		problems[k].Size = p.Size
		problems[k].Page = p.Page
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
func (p *Problem) GetProblemRandom(redisClient *redis.Client) (*mysql.Problems, error) {
	problem, err := mysql.GetProblemRandom()
	if err != nil {
		zap.L().Error("services-GetProblemRandom-GetProblemRandom", zap.Error(err))
		return nil, err
	}

	// 加入redis缓存
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, problem.ProblemID)
	// 将获取的题目列表数据保存到 Redis 缓存中
	encodedData, err := json.Marshal(problem)
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-Marshal ", zap.Error(err))
		return problem, nil
	}

	// 设置缓存的过期时间，你也可以根据具体情况设置适当的缓存时间
	expiration := 5 * time.Hour
	err = redisClient.Set(Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-redisClient.Set ", zap.Error(err))
	}
	return problem, nil
}

// GetProblemListWithCache 获取题目列表，使用 Redis 缓存
func (p *Problem) GetProblemListWithCache(redisClient *redis.Client) (*[]Problem, error) {
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%d:%d", define.GlobalCacheKeyMap.ProblemListPrefix, p.Page, p.Size)
	cachedData, err := redisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var problems []Problem
		err := json.Unmarshal([]byte(cachedData), &problems)
		if err != nil {
			zap.L().Error("services-GetProblemListWithCache-Unmarshal ", zap.Error(err))
			// 从缓存中读取的数据不符合预期的格式，需要从数据库中重新获取
		} else {
			return &problems, nil
		}
	}

	// 缓存中不存在数据，从数据库中获取题目列表
	problems, err := p.GetProblemList()
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-p.GetProblemList ", zap.Error(err))
		return nil, err
	}

	// 将获取的题目列表数据保存到 Redis 缓存中
	encodedData, err := json.Marshal(problems)
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-Marshal ", zap.Error(err))
		return problems, nil
	}

	// 设置缓存的过期时间，你也可以根据具体情况设置适当的缓存时间
	expiration := 5 * time.Hour
	err = redisClient.Set(Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-redisClient.Set ", zap.Error(err))
	}

	return problems, nil
}

// GetProblemDetailWithCache 获取单个题目详细信息
func (p *Problem) GetProblemDetailWithCache(redisClient *redis.Client) (*mysql.Problems, error) {
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, p.ProblemID)

	cachedData, err := redisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var problems mysql.Problems
		err := json.Unmarshal([]byte(cachedData), &problems)
		if err != nil {
			zap.L().Error("services-GetProblemDetailWithCache-Unmarshal ", zap.Error(err))
			// 从缓存中读取的数据不符合预期的格式，需要从数据库中重新获取
		} else {
			return &problems, nil
		}
	}

	// 缓存中不存在数据，从数据库中获取题目列表
	problems, err := p.GetProblemDetail()
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-p.GetProblemDetail ", zap.Error(err))
		return nil, err
	}

	// 将获取的题目列表数据保存到 Redis 缓存中
	encodedData, err := json.Marshal(problems)
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-Marshal ", zap.Error(err))
		return problems, nil
	}

	// 设置缓存的过期时间，你也可以根据具体情况设置适当的缓存时间
	expiration := 5 * time.Hour
	err = redisClient.Set(Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-redisClient.Set ", zap.Error(err))
	}

	return problems, nil
}
