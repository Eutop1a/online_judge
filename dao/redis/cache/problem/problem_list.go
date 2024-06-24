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
	"online_judge/pkg/common_define"
	"time"
)

type CacheProblem struct{}

// GetProblemListWithCache 获取题目列表，使用 Redis 缓存和布隆过滤器
func (p *CacheProblem) GetProblemListWithCache(req request.GetProblemListReq) (response.GetProblemListResp, error) {
	var problems response.GetProblemListResp
	cacheKey := fmt.Sprintf("%s:page-%d:size-%d",
		common_define.GlobalCacheKeyMap.ProblemListPrefix, req.Page, req.Size)
	cachedData, err := redis2.Client.ZRange(redis2.Ctx, cacheKey, 0, -1).Result()

	if err == redis.Nil || len(cachedData) == 0 {
		// 布隆过滤器检查
		if !bloom.ProblemListBloomFilter.TestString(cacheKey) {
			fmt.Println("cache key not exist")
			return problems, common_define.ErrorBloomFilterNotFound
		}

		// 缓存未命中，从数据库获取数据
		problems, err = ProblemService.GetProblemList(req)
		if err != nil {
			zap.L().Error("services-GetProblemListWithCache-GetProblemList", zap.Error(err))
			return problems, err
		}

		// 将数据保存到缓存
		for _, problem := range problems.Data {
			encodedData, err := json.Marshal(problem)
			if err != nil {
				zap.L().Error("services-GetProblemListWithCache-Marshal", zap.Error(err))
				return problems, nil
			}
			redis2.Client.ZAdd(redis2.Ctx, cacheKey, redis.Z{
				Score:  float64(problem.ID),
				Member: encodedData,
			})
		}

		expiration := time.Duration(5+rand.Intn(5)) * time.Hour // 随机过期时间，防止缓存雪崩
		redis2.Client.Expire(redis2.Ctx, cacheKey, expiration)
	} else if err != nil {
		// Redis操作出错
		zap.L().Error("services-GetProblemListWithCache-ZRange", zap.Error(err))
		return problems, err
	} else {
		// 缓存命中，反序列化数据
		problems.Data = make([]*response.ProblemResponse, len(cachedData))
		for i, item := range cachedData {
			var problem response.ProblemResponse
			err = json.Unmarshal([]byte(item), &problem)
			if err != nil {
				zap.L().Error("services-GetProblemListWithCache-Unmarshal", zap.Error(err))
				return problems, err
			}
			problems.Data[i] = &problem
		}
	}

	problems.Count = int64(len(problems.Data))
	problems.Size = req.Size
	problems.Page = req.Page

	return problems, nil
}
