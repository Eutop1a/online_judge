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
	"online_judge/pkg/common_define"
	"strings"
	"time"
)

// SearchProblemWithCache 搜索题目，使用 Redis 缓存和布隆过滤器
func (p *CacheProblem) SearchProblemWithCache(req request.SearchProblemReq) (response.SearchProblemResp, error) {
	var problems response.SearchProblemResp

	keyword := strings.ToLower(req.Title)
	cacheKey := fmt.Sprintf("%s:%s", common_define.GlobalCacheKeyMap.ProblemSearchPrefix, keyword)
	// 使用 HSCAN 进行模糊查询
	cachedData, _, err := redis2.Client.
		HScan(redis2.Ctx, cacheKey, 0, fmt.Sprintf("*%s*", keyword), 10).Result()

	// 缓存未命中
	if err == redis.Nil || len(cachedData) == 0 {
		// 直接去数据库中取
		problem, err := ProblemService.SearchProblem(req)
		if err != nil {
			return problems, err
		}
		// 说明数据库中不存在这样的值，对这个 key 缓存0值
		if len(problem.Data) == 0 {
			// 写入 redis
			err = redis2.Client.HSet(redis2.Ctx, cacheKey, keyword, "").Err()
			if err != nil {
				zap.L().Error("SearchProblemWithCache-HSet", zap.Error(err))
				return problems, err
			}
			return problems, nil
		}
		// 给外部的局部变量赋值
		problems = problem

		// 序列化查询到的数据
		jsonData, err := json.Marshal(problem)
		if err != nil {
			zap.L().Error("SearchProblemWithCache-Marshal", zap.Error(err))
			return problems, err
		}
		// 写入 redis
		err = redis2.Client.HSet(redis2.Ctx, cacheKey, keyword, string(jsonData)).Err()
		if err != nil {
			zap.L().Error("SearchProblemWithCache-HSet", zap.Error(err))
			return problems, err
		}

		// 设置随机的过期时间，防止缓存雪崩
		expiration := time.Duration(5+rand.Intn(5)) * time.Hour
		redis2.Client.Expire(redis2.Ctx, cacheKey, expiration)

	} else if err != nil {
		// Redis操作出错
		zap.L().Error("services-GetProblemListWithCache-ZRange", zap.Error(err))
		return problems, err
	} else {
		// 缓存命中，反序列化数据
		for i := 1; i < len(cachedData); i += 2 {
			if cachedData[i] == "" {
				continue
			}
			var cachedProblem response.SearchProblemResp
			err = json.Unmarshal([]byte(cachedData[i]), &cachedProblem)
			if err != nil {
				zap.L().Error("services-GetProblemListWithCache-Unmarshal", zap.Error(err), zap.String("data", cachedData[i]))
				continue
			}
			problems.Data = append(problems.Data, cachedProblem.Data...)
		}
	}
	p.RecordSearch(keyword)

	return problems, nil
}

//// RecordSearch 记录搜索关键词的热度和最近搜索
//func (p *CacheProblem) RecordSearch(keyword string) {
//	// 记录最热搜索
//	redis2.Client.ZIncrBy(redis2.Ctx, common_define.GlobalCacheKeyMap.HotSearchPrefix, 1, keyword)
//
//	// 记录最近搜索
//	redis2.Client.LPush(redis2.Ctx, common_define.GlobalCacheKeyMap.RecentSearchPrefix, keyword)
//	redis2.Client.LTrim(redis2.Ctx, common_define.GlobalCacheKeyMap.RecentSearchPrefix, 0, 99) // 保留最近100条记录
//}

//// GetHotSearches 获取最热搜索关键词
//func (p *CacheProblem) GetHotSearches(limit int) ([]string, error) {
//	return redis2.Client.ZRevRange(redis2.Ctx, common_define.GlobalCacheKeyMap.HotSearchPrefix, 0, int64(limit-1)).Result()
//}

// GetRecentSearches 获取最近搜索关键词
func (p *CacheProblem) GetRecentSearches(limit int) ([]string, error) {
	return redis2.Client.LRange(redis2.Ctx, common_define.GlobalCacheKeyMap.RecentSearchPrefix, 0, int64(limit-1)).Result()
}

// RecordSearch 记录搜索关键词的热度和最近搜索
func (p *CacheProblem) RecordSearch(keyword string) {
	// 获取当前时间戳
	currentTimestamp := float64(time.Now().Unix())

	// 记录最热搜索，带有时间戳
	luaScript := `
    local keyword = ARGV[1]
    local timestamp = tonumber(ARGV[2])
    local hotSearchKey = KEYS[1]

    -- 获取当前分数
    local currentScore = redis.call("ZSCORE", hotSearchKey, keyword)
    if not currentScore then
        currentScore = 0
    end

    -- 更新分数，使用时间戳和搜索次数结合的方式
    local newScore = tonumber(currentScore) + 1 + timestamp * 0.0001
    redis.call("ZADD", hotSearchKey, newScore, keyword)
    return 1
    `
	script := redis.NewScript(luaScript)
	_, err := script.Run(redis2.Ctx, redis2.Client, []string{common_define.GlobalCacheKeyMap.HotSearchPrefix}, keyword, currentTimestamp).Result()
	if err != nil {
		zap.L().Error("RecordSearch-HotSearch-LuaScript", zap.Error(err))
	}

	// 记录最近搜索，去重
	luaScriptRecent := `
    local keyword = ARGV[1]
    local recentSearchKey = KEYS[1]

    -- 移除列表中的现有关键词
    redis.call("LREM", recentSearchKey, 0, keyword)
    -- 将关键词推到列表最前面
    redis.call("LPUSH", recentSearchKey, keyword)
    -- 保留列表的前100个元素
    redis.call("LTRIM", recentSearchKey, 0, 99)
    return 1
    `
	scriptRecent := redis.NewScript(luaScriptRecent)
	_, err = scriptRecent.Run(redis2.Ctx, redis2.Client, []string{common_define.GlobalCacheKeyMap.RecentSearchPrefix}, keyword).Result()
	if err != nil {
		zap.L().Error("RecordSearch-RecentSearch-LuaScript", zap.Error(err))
	}
}

// GetHotSearches 获取最热搜索关键词
func (p *CacheProblem) GetHotSearches(limit int) ([]string, error) {
	hotSearchKey := common_define.GlobalCacheKeyMap.HotSearchPrefix

	luaScript := `
    local hotSearchKey = KEYS[1]
    local limit = tonumber(ARGV[1])
    local currentTimestamp = tonumber(ARGV[2])

    -- 获取所有带有分数的成员
    local members = redis.call("ZRANGE", hotSearchKey, 0, -1, "WITHSCORES")
    local adjustedMembers = {}

    for i = 1, #members, 2 do
        local keyword = members[i]
        local score = tonumber(members[i + 1])
        local adjustedScore = score - (currentTimestamp - math.floor(score * 10000) / 10000) * 0.0001
        table.insert(adjustedMembers, {keyword, adjustedScore})
    end

    -- 按照调整后的分数排序
    table.sort(adjustedMembers, function(a, b) return a[2] > b[2] end)

    -- 取前limit个
    local result = {}
    for i = 1, limit do
        if adjustedMembers[i] then
            table.insert(result, adjustedMembers[i][1])
        end
    end

    return result
    `
	script := redis.NewScript(luaScript)
	currentTimestamp := float64(time.Now().Unix())
	res, err := script.Run(redis2.Ctx, redis2.Client, []string{hotSearchKey}, limit, currentTimestamp).Result()
	if err != nil {
		zap.L().Error("GetHotSearches-LuaScript", zap.Error(err))
		return nil, err
	}

	// 转换结果类型
	var result []string
	for _, r := range res.([]interface{}) {
		result = append(result, r.(string))
	}

	return result, nil
}
