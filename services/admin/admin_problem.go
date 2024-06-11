package admin

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"online_judge/consts"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"online_judge/pkg/define"
	"os"
	"path/filepath"
)

type AdminProblemService struct{}

// CreateProblem 创建题目
func (p *AdminProblemService) CreateProblem(request request.AdminCreateProblemReq) (response response.Response) {
	// 检查题目标题是否重复
	exists, err := mysql.CheckProblemTitleExists(request.Title)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-CreateProblem-CheckProblemTitle ", zap.Error(err))
		return
	}
	if exists {
		response.Code = resp_code.ProblemAlreadyExist
		zap.L().Error("services-CreateProblem-CheckProblemTitle " +
			fmt.Sprintf("title %s aleardy exist", request.Title))
		return
	}

	// 创建题目
	err = mysql.CreateProblem(&mysql.Problems{
		ProblemID:  request.ProblemID,
		Title:      request.Title,
		Content:    request.Content,
		Difficulty: request.Difficulty,
		MaxRuntime: request.MaxRuntime,
		MaxMemory:  request.MaxMemory,
		TestCases:  p.convertTestCases(request.TestCases),
	})
	if err != nil {
		response.Code = resp_code.CreateProblemError
		zap.L().Error("services-CreateProblem-CreateProblem ", zap.Error(err))
		return
	}
	// 添加成功后删除缓存
	if err := p.deleteCacheByPrefix(request.RedisClient, define.GlobalCacheKeyMap.ProblemListPrefix); err != nil {
		zap.L().Error("services-CreateProblem-deleteCacheByPrefix ", zap.Error(err))
		response.Code = resp_code.DeleteCacheError
		return
	}
	// 删除特定问题的缓存（如果存在）
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, request.ProblemID)
	if err := request.RedisClient.Del(request.Ctx, cacheKey).Err(); err != nil {
		zap.L().Error("services-CreateProblem-redisClient.Del ", zap.Error(err))
		response.Code = resp_code.DeleteCacheError
		return
	}
	response.Code = resp_code.Success
	return
}

// UpdateProblem 更新题目
func (p *AdminProblemService) UpdateProblem(request request.AdminUpdateProblemReq) (response response.Response) {
	// 检查题目id是否存在
	exists, err := mysql.CheckProblemIDExists(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateProblem-CheckProblemID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.ProblemNotExist
		zap.L().Error("services-UpdateProblem-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", request.ProblemID))
		return
	}
	if request.Title != "" {
		// 检查题目标题是否重复
		exists, err = mysql.CheckProblemTitleExists(request.Title)
		if err != nil {
			response.Code = resp_code.SearchDBError
			zap.L().Error("services-CreateProblem-CheckProblemTitle ", zap.Error(err))
			return
		}
		if exists {
			response.Code = resp_code.ProblemAlreadyExist
			zap.L().Error("services-CreateProblem-CheckProblemTitle " +
				fmt.Sprintf("title %s aleardy exist", request.Title))
			return
		}
	}

	err = mysql.UpdateProblem(&mysql.Problems{
		ProblemID:  request.ProblemID,
		Title:      request.Title,
		Content:    request.Content,
		Difficulty: request.Difficulty,
		MaxRuntime: request.MaxRuntime,
		MaxMemory:  request.MaxMemory,
		TestCases:  p.convertTestCases(request.TestCases),
	})

	if err != nil {
		zap.L().Error("services-UpdateProblem-UpdateProblem ", zap.Error(err))
		response.Code = resp_code.InternalServerError
		return
	}
	// 更新成功后删除缓存
	if err := p.deleteCacheByPrefix(request.RedisClient, define.GlobalCacheKeyMap.ProblemListPrefix); err != nil {
		zap.L().Error("services-CreateProblem-deleteCacheByPrefix ", zap.Error(err))
		response.Code = resp_code.DeleteCacheError
		return
	}
	// 删除特定问题的缓存（如果存在）
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, request.ProblemID)
	if err := request.RedisClient.Del(request.Ctx, cacheKey).Err(); err != nil {
		zap.L().Error("services-CreateProblem-redisClient.Del ", zap.Error(err))
		response.Code = resp_code.DeleteCacheError
		return
	}
	response.Code = resp_code.Success
	return
}

// DeleteProblem 删除题目
func (p *AdminProblemService) DeleteProblem(request request.AdminDeleteProblemReq) (response response.Response) {
	// 检查题目id是否存在
	exists, err := mysql.CheckProblemIDExists(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateProblem-CheckProblemID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.ProblemNotExist
		zap.L().Error("services-UpdateProblem-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", request.ProblemID))
		return
	}

	// 删除题目
	err = mysql.DeleteProblem(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteProblem-DeleteProblem  ", zap.Error(err))
		return
	}
	// 删除成功后删除缓存
	if err := p.deleteCacheByPrefix(request.RedisClient, define.GlobalCacheKeyMap.ProblemListPrefix); err != nil {
		zap.L().Error("services-CreateProblem-deleteCacheByPrefix ", zap.Error(err))
		response.Code = resp_code.DeleteCacheError
		return
	}
	// 删除特定问题的缓存（如果存在）
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, request.ProblemID)
	if err := request.RedisClient.Del(request.Ctx, cacheKey).Err(); err != nil {
		zap.L().Error("services-CreateProblem-redisClient.Del ", zap.Error(err))
		response.Code = resp_code.DeleteCacheError
		return
	}
	response.Code = resp_code.Success
	return
}

// CreateProblemWithFile 创建测试用例为file的题目
func (p *AdminProblemService) CreateProblemWithFile(request request.AdminCreateProblemWithFileReq) (response response.Response) {
	// 检查题目标题是否重复
	var problemNum int64
	err := mysql.CheckProblemIDWithFile(request.Title, &problemNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-CreateProblemWithFile-CheckProblemTitle ", zap.Error(err))
		return
	case problemNum > 0: // 题目已经存在
		response.Code = resp_code.ProblemAlreadyExist
		zap.L().Error("services-CreateProblemWithFile-CheckProblemTitle " +
			fmt.Sprintf("title %s aleardy exist", request.Title))
		return
	}
	// 创建题目
	err = mysql.CreateProblemWithFile(&mysql.ProblemWithFile{
		ProblemID:    request.ProblemID,
		Title:        request.Title,
		Content:      request.Content,
		Difficulty:   request.Difficulty,
		MaxRuntime:   request.MaxRuntime,
		MaxMemory:    request.MaxMemory,
		InputPath:    request.InputDst,
		ExpectedPath: request.ExpectedDst,
	})

	if err != nil {
		response.Code = resp_code.CreateProblemError
		zap.L().Error("services-CreateProblemWithFile-CreateProblemWithFile ", zap.Error(err))
		return
	}
	response.Code = resp_code.Success
	return
}

// DeleteProblemWithFile 删除题目
func (p *AdminProblemService) DeleteProblemWithFile(request request.AdminDeleteProblemWithFileReq) (response response.Response) {
	// 检查题目id是否存在
	var idNum int64
	err := mysql.CheckProblemIDWithFile(request.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteProblemWithFile-CheckProblemIDWithFile ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = resp_code.ProblemNotExist
		zap.L().Error("services-DeleteProblemWithFile-CheckProblemIDWithFile " +
			fmt.Sprintf("problemID %s do not exist", request.ProblemID))
		return
	}
	// 删除题目
	problemID, err := mysql.DeleteProblemWithFile(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteProblemWithFile-DeleteProblemWithFile ", zap.Error(err))
		return
	}
	//fmt.Println("path: ", filepath.Join(consts.FilePath, problemID))
	err = os.RemoveAll(filepath.Join(consts.FilePath, problemID))

	if err != nil {
		response.Code = resp_code.RemoveTestFileError
		zap.L().Error("services-DeleteProblemWithFile-Remove ", zap.Error(err))
		return
	}
	response.Code = resp_code.Success
	return
}

// DeleteProblemTestCaseWithFile 删除题目测试用例文件
func (p *AdminProblemService) DeleteProblemTestCaseWithFile(request request.AdminUpdateProblemWithFileReq) (response response.Response) {
	// 检查题目id是否存在
	var idNum int64
	err := mysql.CheckProblemIDWithFile(request.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteProblemTestCaseWithFile-CheckProblemID ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = resp_code.ProblemNotExist
		zap.L().Error("services-DeleteProblemTestCaseWithFile-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", request.ProblemID))
		return
	}

	// 检查题目标题是否存在
	var titleNum int64
	err = mysql.CheckProblemTitleWithFile(request.Title, &titleNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateProblemWithFile-CheckProblemTitle", zap.Error(err))
		return
	case titleNum != 0: // 题目title已存在
		response.Code = resp_code.ProblemAlreadyExist
		zap.L().Error("services-UpdateProblemWithFile-CheckProblemTitle" +
			fmt.Sprintf("problem title %s already exist", request.Title))
		return
	}
	err = os.RemoveAll(filepath.Join(consts.FilePath, request.ProblemID))

	if err != nil {
		response.Code = resp_code.RemoveTestFileError
		zap.L().Error("services-DeleteProblemTestCaseWithFile-Remove ", zap.Error(err))
		return
	}
	response.Code = resp_code.Success
	return
}

// UpdateProblemWithFile 更新题目
func (p *AdminProblemService) UpdateProblemWithFile(request request.AdminUpdateProblemWithFileReq) (response response.Response) {

	err := mysql.UpdateProblemWithFile(&mysql.ProblemWithFile{
		ProblemID:    request.ProblemID,
		Title:        request.Title,
		Content:      request.Content,
		Difficulty:   request.Difficulty,
		MaxRuntime:   request.MaxRuntime,
		MaxMemory:    request.MaxMemory,
		InputPath:    request.InputDst,
		ExpectedPath: request.ExpectedDst,
	})
	if err != nil {
		zap.L().Error("services-UpdateProblemWithFile-UpdateProblemWithFile ", zap.Error(err))
		response.Code = resp_code.InternalServerError
		return
	}
	response.Code = resp_code.Success
	return
}

func (p *AdminProblemService) convertTestCases(testCases []*request.TestCase) []*mysql.TestCase {
	// 提前转换类型
	var convertedTestCases []*mysql.TestCase
	for _, tc := range testCases {
		// 进行类型转换
		convertedTestCases = append(convertedTestCases, &mysql.TestCase{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		})
	}
	return convertedTestCases
}

func (p *AdminProblemService) deleteCacheByPrefix(redisClient *redis.Client, prefix string) error {
	ctx := context.Background()
	iter := redisClient.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
