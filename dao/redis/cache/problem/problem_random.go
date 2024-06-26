package problem

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	redis2 "online_judge/dao/redis"
	"online_judge/models/problem/request"
	"online_judge/models/problem/response"
	"online_judge/pkg/define"
	"time"
)

// ProblemIDListCacheGet 从 Cache 中随机获取题目 ID
func (p *CacheProblem) ProblemIDListCacheGet() (string, error) {
	count, err := redis2.Client.ZCard(redis2.Ctx, define.GlobalCacheKeyMap.ProblemListPrefix).Result()
	if err != nil {
		zap.L().Error("redis2.Client.ZCARD fail", zap.Error(err))
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	randomIdx := rand.Intn(int(count))

	randomProblemID, err := redis2.Client.ZRange(redis2.Ctx,
		define.GlobalCacheKeyMap.ProblemListPrefix, int64(randomIdx), int64(randomIdx)).Result()
	if err != nil {
		zap.L().Error("redis2.Client.ZRANGE fail", zap.Error(err))
		return "", err
	}
	if len(randomProblemID) == 0 {
		return "", nil
	}
	return randomProblemID[0], nil
}

// ProblemIDListCacheRandom 根据题目ID 获取题目
func (p *CacheProblem) ProblemIDListCacheRandom() (response.ProblemDetailResponse, error) {
	var problem response.ProblemDetailResponse
	randomProblemID, err := p.ProblemIDListCacheGet()
	if err != nil {
		return problem, err
	}
	if randomProblemID == "" {
		return problem, define.ErrNoProblemIDFound
	}

	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemListPrefix, randomProblemID)
	cachedData, err := redis2.Client.Get(redis2.Ctx, cacheKey).Result()

	if err == redis.Nil || cachedData == "" {
		problemDetail, err := ProblemService.GetProblemDetail(request.GetProblemDetailReq{ProblemID: randomProblemID})
		if err != nil {
			zap.L().Error("ProblemService.GetProblemDetail fail", zap.Error(err))
			return problem, err
		}

		encodeData, err := json.Marshal(problemDetail)
		if err != nil {
			zap.L().Error("ProblemDetailResponse.Marshal fail", zap.Error(err))
			return problem, err
		}

		expiration := time.Duration(5+rand.Intn(5)) * time.Hour
		_, err = redis2.Client.TxPipelined(redis2.Ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(redis2.Ctx, cacheKey, string(encodeData), expiration)
			return nil
		})
		if err != nil {
			zap.L().Error("redis2.Client.Set fail", zap.Error(err))
			return problem, err
		}

		problem = *problemDetail
	} else if err != nil {
		zap.L().Error("redis2.Client.Get fail", zap.Error(err))
		return problem, err
	} else {
		err = json.Unmarshal([]byte(cachedData), &problem)
		if err != nil {
			zap.L().Error("ProblemDetailResponse.Unmarshal fail", zap.Error(err))
			return problem, err
		}
	}

	return problem, nil
}
