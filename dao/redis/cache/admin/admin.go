package admin

import (
	"fmt"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/redis"
	"online_judge/dao/redis/bloom"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"online_judge/pkg/common_define"
)

// CreateProblem 创建题目
func (p *CacheGroup) CreateProblem(request request.AdminCreateProblemReq) response.Response {
	resp := AdminService.CreateProblem(request)

	if resp.Code != resp_code.Success {
		return resp
	}
	// 刷新布隆过滤器
	bloom.ReBuildBloomFilters()

	// 删除题目列表的缓存
	if err := p.DeleteProblemListCacheByPrefix(common_define.GlobalCacheKeyMap.ProblemListPrefix); err != nil {
		zap.L().Error("services-UpdateProblem-DeleteProblemListCacheByPrefix ", zap.Error(err))
		resp.Code = resp_code.DeleteCacheError
		return resp
	}
	resp.Code = resp_code.Success
	return resp
}

// UpdateProblem 更新题目
func (p *CacheGroup) UpdateProblem(request request.AdminUpdateProblemReq) response.Response {
	resp := AdminService.UpdateProblem(request)

	if resp.Code != resp_code.Success {
		return resp
	}
	// 刷新布隆过滤器
	bloom.ReBuildBloomFilters()

	// 删除题目列表的缓存
	if err := p.DeleteProblemListCacheByPrefix(common_define.GlobalCacheKeyMap.ProblemListPrefix); err != nil {
		zap.L().Error("services-UpdateProblem-DeleteProblemListCacheByPrefix ", zap.Error(err))
		resp.Code = resp_code.DeleteCacheError
		return resp
	}

	// 删除特定问题的缓存（如果存在）
	cacheKey := fmt.Sprintf("%s:%s", common_define.GlobalCacheKeyMap.ProblemDetailPrefix, request.ProblemID)
	if err := p.DeleteProblemDetailCacheByPrefix(cacheKey); err != nil {
		zap.L().Error("services-UpdateProblem-DeleteProblemDetailCacheByPrefix ", zap.Error(err))
		resp.Code = resp_code.DeleteCacheError
		return resp
	}
	resp.Code = resp_code.Success
	return resp
}

func (p *CacheGroup) DeleteProblem(request request.AdminDeleteProblemReq) response.Response {
	resp := AdminService.DeleteProblem(request)
	if resp.Code != resp_code.Success {
		return resp
	}
	// 刷新布隆过滤器
	bloom.ReBuildBloomFilters()

	// 删除题目列表的缓存
	if err := p.DeleteProblemListCacheByPrefix(common_define.GlobalCacheKeyMap.ProblemListPrefix); err != nil {
		zap.L().Error("services-UpdateProblem-DeleteProblemListCacheByPrefix ", zap.Error(err))
		resp.Code = resp_code.DeleteCacheError
		return resp
	}
	// 删除特定问题的缓存（如果存在）
	cacheKey := fmt.Sprintf("%s:%s", common_define.GlobalCacheKeyMap.ProblemDetailPrefix, request.ProblemID)
	if err := p.DeleteProblemDetailCacheByPrefix(cacheKey); err != nil {
		zap.L().Error("services-UpdateProblem-DeleteProblemDetailCacheByPrefix ", zap.Error(err))
		resp.Code = resp_code.DeleteCacheError
		return resp
	}
	resp.Code = resp_code.Success
	return resp
}

// DeleteProblemListCacheByPrefix 删除问题列表的缓存
func (p *CacheGroup) DeleteProblemListCacheByPrefix(prefix string) error {
	iter := redis.Client.Scan(redis.Ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(redis.Ctx) {
		if err := redis.Client.Del(redis.Ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

// DeleteProblemDetailCacheByPrefix 删除问题详细信息的缓存
func (p *CacheGroup) DeleteProblemDetailCacheByPrefix(cacheKey string) error {
	if err := redis.Client.Del(redis.Ctx, cacheKey).Err(); err != nil {
		return err
	}
	return nil
}
