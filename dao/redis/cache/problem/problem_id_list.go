package problem

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	"online_judge/dao/mysql"
	redis2 "online_judge/dao/redis"
	"online_judge/pkg/define"
)

// ProblemIDListCacheInit 获取所有的 problemID，并存储到 redis 中
func ProblemIDListCacheInit() {
	data, err := mysql.GetAllProblem()
	if err != nil {
		zap.L().Error("mysql.GetAllProblem fail", zap.Error(err))
		return
	}
	// 将 problemID 加入有序集合中，分数为随机值
	for _, problem := range data {
		score := rand.Float64()
		redis2.Client.ZAdd(redis2.Ctx, define.GlobalCacheKeyMap.ProblemListPrefix, redis.Z{
			Score:  score,
			Member: problem.ProblemID,
		})
	}
}
