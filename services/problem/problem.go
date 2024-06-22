package problem

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	problems.Data = make([]*response.ProblemResponse, len(problemList))
	for k, v := range problemList {
		problems.Data[k] = &response.ProblemResponse{ // 为每个元素分配内存
			ID:         v.ID,
			ProblemID:  v.ProblemID,
			Title:      v.Title,
			Difficulty: v.Difficulty,
			Categories: make([]response.CategoryResponse, len(v.ProblemCategories)),
		}
		// 将 ProblemCategories 复制到响应中
		for i, pc := range v.ProblemCategories {
			if pc.Category != nil {
				problems.Data[k].Categories[i] = response.CategoryResponse{
					CategoryID: pc.CategoryIdentity,
					Name:       pc.Category.Name,
					ParentName: pc.Category.ParentName,
				}
			}
		}
	}
	problems.Count = count
	problems.Size = req.Size
	problems.Page = req.Page

	return
}

// GetProblemDetail 获取单个题目详细信息
func (p *ProblemService) GetProblemDetail(req request.GetProblemDetailReq) (*response.ProblemResponse, error) {
	// 从数据库中获取题目详细，包括测试用例和分类
	problem, err := mysql.GetProblemDetail(req.ProblemID)
	if err != nil {
		zap.L().Error("services-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return nil, err
	}

	// 构建返回的结构体
	problemResp := &response.ProblemResponse{
		ID:         problem.ID,
		ProblemID:  problem.ProblemID,
		Title:      problem.Title,
		Difficulty: problem.Difficulty,
		Categories: make([]response.CategoryResponse, len(problem.ProblemCategories)),
		TestCases:  make([]response.TestCaseResponse, len(problem.TestCases)),
	}

	// 赋值分类
	for i, pc := range problem.ProblemCategories {
		if pc.Category != nil {
			problemResp.Categories[i] = response.CategoryResponse{
				CategoryID: pc.CategoryIdentity,
				Name:       pc.Category.Name,
				ParentName: pc.Category.ParentName,
			}
		}
	}

	// 赋值
	for i, tc := range problem.TestCases {
		problemResp.TestCases[i] = response.TestCaseResponse{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}

	return problemResp, nil
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
func (p *ProblemService) GetProblemListWithCache(req request.GetProblemListReq) (response.GetProblemListResp, error) {
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%d:%d", define.GlobalCacheKeyMap.ProblemListPrefix, req.Page, req.Size)
	cachedData, err := req.RedisClient.Get(req.Ctx, cacheKey).Result()
	if err == nil {
		var problems response.GetProblemListResp
		err = json.Unmarshal([]byte(cachedData), &problems)
		if err != nil {
			zap.L().Error("services-GetProblemListWithCache-Unmarshal ", zap.Error(err))
			// 从缓存中读取的数据不符合预期的格式，需要从数据库中重新获取
		} else {
			return problems, nil
		}
	}

	// 缓存中不存在数据，从数据库中获取题目列表
	problems, err := p.GetProblemList(req)
	if err != nil {
		zap.L().Error("services-GetProblemListWithCache-p.GetProblemList ", zap.Error(err))
		return problems, err
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
func (p *ProblemService) GetProblemDetailWithCache(req request.GetProblemDetailReq) (*response.ProblemResponse, error) {
	// 尝试从缓存中获取题目列表
	cacheKey := fmt.Sprintf("%s:%s", define.GlobalCacheKeyMap.ProblemDetailPrefix, req.ProblemID)

	cachedData, err := req.RedisClient.Get(req.Ctx, cacheKey).Result()
	if err == nil {
		var problem response.ProblemResponse
		err := json.Unmarshal([]byte(cachedData), &problem)
		if err != nil {
			zap.L().Error("services-GetProblemDetailWithCache-Unmarshal ", zap.Error(err))
			// 从缓存中读取的数据不符合预期的格式，需要从数据库中重新获取
		} else {
			return &problem, nil
		}
	}

	// 缓存中不存在数据，从数据库中获取题目列表
	problem, err := p.GetProblemDetail(req)
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-p.GetProblemDetail ", zap.Error(err))
		return nil, err
	}

	// 将获取的题目列表数据保存到 Redis 缓存中
	encodedData, err := json.Marshal(problem)
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-Marshal ", zap.Error(err))
		return problem, nil
	}

	// 设置缓存的过期时间，你也可以根据具体情况设置适当的缓存时间
	expiration := 5 * time.Hour
	err = req.RedisClient.Set(req.Ctx, cacheKey, encodedData, expiration).Err()
	if err != nil {
		zap.L().Error("services-GetProblemDetailWithCache-redisClient.Set ", zap.Error(err))
	}

	return problem, nil
}

func (p *ProblemService) SearchProblem(req request.SearchProblemReq) (response.SearchProblemResp, error) {
	var problems []mysql.Problems
	searchQuery := "%" + req.Msg + "%"

	var data response.SearchProblemResp
	// 根据提供的信息搜索，若没有记录则返回空，若出错返回500
	err := mysql.SearchProblemByMsg(&problems, searchQuery)
	if err != nil {
		// 没找到对应的记录
		if err == gorm.ErrRecordNotFound {
			zap.L().Error("services-SearchProblem-SearchProblemByMsg ", zap.Error(err))
			return data, err
		}
		// 数据库内部错误
		zap.L().Error("services-SearchProblem-SearchProblemByMsg ", zap.Error(err))
		return data, err
	}

	data.Data = make([]*response.ProblemResponse, len(problems))
	for i, problem := range problems {
		data.Data[i] = &response.ProblemResponse{
			ID:         problem.ID,
			ProblemID:  problem.ProblemID,
			Title:      problem.Title,
			Difficulty: problem.Difficulty,
			Categories: make([]response.CategoryResponse, len(problem.ProblemCategories)),
		}
		for j, pc := range problem.ProblemCategories {
			if pc.Category != nil {
				data.Data[i].Categories[j] = response.CategoryResponse{
					CategoryID: pc.CategoryIdentity,
					Name:       pc.Category.Name,
					ParentName: pc.Category.ParentName,
				}
			}
		}
	}
	return data, nil
}
func (p *ProblemService) _GetProblemListByCategory(categoryName string) (*response.GetProblemListResp, error) {
	// 先在 mysql 里面检测这个 categoryName 存在不存在
	ok, err := mysql.CheckCategoryByName(categoryName)
	if err != nil {
		zap.L().Error("services-GetProblemListByCategory CheckCategoryByName", zap.Error(err))
		return nil, err
	}
	if !ok {
		zap.L().Error("services-GetProblemListByCategory categoryName not exist")
		return nil, nil
	}

	// 获取 categoryName 对应的 categoryID
	categoryID, err := mysql.GetCategoryID(categoryName)
	if err != nil {
		zap.L().Error("services-GetProblemListByCategory GetCategoryID ", zap.Error(err))
		return nil, err
	}

	problems, err := mysql.GetProblemListByCategory(categoryID)
	if err != nil {
		zap.L().Error("services-GetProblemListByCategory-GetProblemListByCategory ", zap.Error(err))
		return nil, err
	}

	var resp response.GetProblemListResp
	resp.Data = make([]*response.ProblemResponse, len(problems))
	for i, problem := range problems {
		resp.Data[i] = &response.ProblemResponse{
			ID:         problem.ID,
			ProblemID:  problem.ProblemID,
			Title:      problem.Title,
			Difficulty: problem.Difficulty,
			Categories: make([]response.CategoryResponse, len(problem.ProblemCategories)),
		}
		for j, pc := range problem.ProblemCategories {
			if pc.Category != nil {
				resp.Data[i].Categories[j] = response.CategoryResponse{
					CategoryID: pc.CategoryIdentity,
					Name:       pc.Category.Name,
					ParentName: pc.Category.ParentName,
				}
			}
		}
	}
	return &resp, nil
}

func (p *ProblemService) GetProblemListByCategory(categoryName string) (*response.GetProblemListResp, error) {
	problems, err := mysql.GetProblemsByCategoryName(categoryName)
	if err != nil {
		zap.L().Error("services-GetProblemListByCategory GetProblemsByCategoryName", zap.Error(err))
		return nil, err
	}
	if len(problems) == 0 {
		zap.L().Error("services-GetProblemListByCategory categoryName not exist or no problems found")
		return nil, nil
	}

	var resp response.GetProblemListResp
	resp.Data = make([]*response.ProblemResponse, len(problems))
	for i, problem := range problems {
		resp.Data[i] = &response.ProblemResponse{
			ID:         problem.ID,
			ProblemID:  problem.ProblemID,
			Title:      problem.Title,
			Difficulty: problem.Difficulty,
			Categories: make([]response.CategoryResponse, len(problem.ProblemCategories)),
		}
		for j, pc := range problem.ProblemCategories {
			if pc.Category != nil {
				resp.Data[i].Categories[j] = response.CategoryResponse{
					CategoryID: pc.CategoryIdentity,
					Name:       pc.Category.Name,
					ParentName: pc.Category.ParentName,
				}
			}
		}
	}
	return &resp, nil
}

// GetCategoryList 获取分类列表
func (p *ProblemService) GetCategoryList() ([]*response.CategoryResponse, error) {
	categoryList, err := mysql.GetCategoryList()
	if err != nil {
		zap.L().Error("services-GetCategoryList-GetCategoryList ", zap.Error(err))
		return nil, err
	}
	var resp []*response.CategoryResponse
	for _, category := range categoryList {
		resp = append(resp, &response.CategoryResponse{
			Name:       category.Name,
			CategoryID: category.CategoryID,
			ParentName: category.ParentName,
		})
	}
	return resp, nil
}
