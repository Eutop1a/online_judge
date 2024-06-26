package problem

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	redis2 "online_judge/dao/redis"
	"online_judge/dao/redis/bloom"
	"online_judge/models/problem/request"
	"online_judge/models/problem/response"
	"online_judge/pkg/define"
	"time"
)

// GetProblemDetailWithCache 获取单个题目详细信息
func (p *CacheProblem) GetProblemDetailWithCache(req request.GetProblemDetailReq) (*response.ProblemDetailResponse, error) {
	var problem *response.ProblemDetailResponse
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, req.ProblemID)
	//cachedData, err := redis2.Client.Get(redis2.Ctx, cacheKey).Result()
	cachedData, err := redis2.Client.HGet(redis2.Ctx, cacheKey, req.ProblemID).Result()

	if err == redis.Nil || len(cachedData) == 0 {
		// 布隆过滤器检查
		if !bloom.ProblemDetailBloomFilter.TestString(cacheKey) {
			zap.L().Error("cache key not exist",
				zap.String("problem_id", req.ProblemID),
			)
			return nil, define.ErrProblemIDNotFound
		}

		// 缓存未命中，从数据库中获取题目详细信息
		problem, err = ProblemService.GetProblemDetail(req)
		if err != nil {
			zap.L().Error("services-GetProblemDetailWithCache-GetProblemDetail ",
				zap.Error(err),
			)
			return nil, err
		}

		// 将数据保存到缓存
		encodeData, err := json.Marshal(problem)
		if err != nil {
			zap.L().Error("services-GetProblemDetailWithCache-Marshal", zap.Error(err))
			return nil, err
		}
		err = redis2.Client.HSet(redis2.Ctx, cacheKey, req.ProblemID, string(encodeData)).Err()
		if err != nil {
			zap.L().Error("services-GetProblemDetailWithCache-HSet", zap.Error(err))
		}

		// 设置随机的过期时间，防止缓存雪崩
		expiration := time.Duration(5+rand.Intn(5)) * time.Hour

		redis2.Client.Expire(redis2.Ctx, cacheKey, expiration)
	} else if err != nil {
		// Redis 操作出错
		zap.L().Error("services-GetProblemDetailWithCache-Set", zap.Error(err))
		return problem, err
	} else {
		// 缓存命中，反序列化数据
		err = json.Unmarshal([]byte(cachedData), &problem)
		if err != nil {
			zap.L().Error("services-GetProblemDetailWithCache-Unmarshal", zap.Error(err))
			return nil, err
		}
	}
	return problem, nil
}
