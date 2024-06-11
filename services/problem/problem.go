package problem

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"online_judge/dao/mysql"
	"online_judge/models/problem/request"
	"online_judge/models/problem/response"
	"online_judge/pkg/define"
	"time"
)

type ProblemService struct{}

// GetProblemList 获取题目列表
func (p *ProblemService) GetProblemList(req request.GetProblemListReq) (problems response.GetProblemListResp, err error) {
	var count int64
	problemList, err := mysql.GetProblemList(req.Page, req.Size, &count)
	if err != nil {
		zap.L().Error("services-GetProblemList-GetProblemList ", zap.Error(err))
		return
	}

	problems.Data = make([]*mysql.Problems, len(problemList))
	for k, v := range problemList {
		problems.Data[k] = &mysql.Problems{ // 为每个元素分配内存
			ID:         v.ID,
			ProblemID:  v.ProblemID,
			Title:      v.Title,
			Difficulty: v.Difficulty,
		}
	}
	problems.Count = count
	problems.Size = req.Size
	problems.Page = req.Page

	return
}

// GetProblemDetail 获取单个题目详细信息
func (p *ProblemService) GetProblemDetail(req request.GetProblemDetailReq) (*mysql.Problems, error) {
	data, err := mysql.GetProblemDetail(req.ProblemID)
	if err != nil {
		zap.L().Error("services-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return nil, err
	}
	return data, nil
}

// GetProblemID 获取题目ID
func (p *ProblemService) GetProblemID(req request.GetProblemIDReq) (problemID string, err error) {
	problemID, err = mysql.GetProblemID(req.Title)
	if err != nil {
		zap.L().Error("services-GetProblemID-GetProblemID", zap.Error(err))
		return "", err
	}
	return
}

// GetProblemRandom 随机获取一个题目
func (p *ProblemService) GetProblemRandom(req request.GetProblemRandomReq) (*mysql.Problems, error) {
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
	err = req.RedisClient.Set(req.Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-redisClient.Set ", zap.Error(err))
	}
	return problem, nil
}

// GetProblemListWithCache 获取题目列表，使用 Redis 缓存
func (p *ProblemService) GetProblemListWithCache(req request.GetProblemListReq) (problems response.GetProblemListResp, err error) {
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%d:%d", define.GlobalCacheKeyMap.ProblemListPrefix, req.Page, req.Size)
	cachedData, err := req.RedisClient.Get(req.Ctx, cacheKey).Result()
	if err == nil {
		var problem []*mysql.Problems
		err = json.Unmarshal([]byte(cachedData), &problem)
		if err != nil {
			zap.L().Error("services-GetProblemListWithCache-Unmarshal ", zap.Error(err))
			// 从缓存中读取的数据不符合预期的格式，需要从数据库中重新获取
		} else {
			problems.Data = problem
			return
		}
	}

	// 缓存中不存在数据，从数据库中获取题目列表
	problems, err = p.GetProblemList(req)
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-p.GetProblemList ", zap.Error(err))
		return
	}

	// 将获取的题目列表数据保存到 Redis 缓存中
	encodedData, err := json.Marshal(problems)
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-Marshal ", zap.Error(err))
		return problems, nil
	}

	// 设置缓存的过期时间，你也可以根据具体情况设置适当的缓存时间
	expiration := 5 * time.Hour
	err = req.RedisClient.Set(req.Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-redisClient.Set ", zap.Error(err))
	}

	return problems, nil
}

// GetProblemDetailWithCache 获取单个题目详细信息
func (p *ProblemService) GetProblemDetailWithCache(req request.GetProblemDetailReq) (*mysql.Problems, error) {
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, req.ProblemID)

	cachedData, err := req.RedisClient.Get(req.Ctx, cacheKey).Result()
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
	problems, err := p.GetProblemDetail(req)
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
	err = req.RedisClient.Set(req.Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-redisClient.Set ", zap.Error(err))
	}

	return problems, nil
}
