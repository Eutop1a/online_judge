package problem

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"online_judge/models/common/response"
	"online_judge/models/problem/request"
	"online_judge/pkg/define"
	"strconv"
)

type ApiProblem struct{}

// GetProblemList 获取题目列表接口
// @Tags Problem API
// @Summary 获取题目列表
// @Description 获取题目列表接口
// @Param page query int false "input current page num, default: 1"
// @Param size query int false "pageSize, default: 10"
// @Success 200 {object} common.GetProblemListResponse "获取题目列表成功"
// @Failure 200 {object} common.GetProblemListResponse "需要登录"
// @Failure 200 {object} common.GetProblemListResponse "服务器内部错误"
// @Router /problem/list [GET]
func (p *ApiProblem) GetProblemList(c *gin.Context) {
	var req request.GetProblemListReq

	req.Size, _ = strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	req.Page, _ = strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))

	data, err := ProblemCache.GetProblemListWithCache(req)
	//data, err := ProblemService.GetProblemListWithCache(req)
	//data, err := getProblemList.GetProblemList()
	if err != nil {
		if err == define.ErrBloomFilterNotFound {
			response.ResponseError(c, response.CodeProblemListNotFound)
			zap.L().Error("controller-GetProblemList-GetProblemList ", zap.Error(err))
			return
		}
		response.ResponseError(c, response.CodeInternalServerError)
		zap.L().Error("controller-GetProblemList-GetProblemList ", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, data)
}

// GetProblemDetail 获取单个题目详细接口
// @Tags Problem API
// @Summary 获取单个题目详细
// @Description 获取单个题目详细接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// // @Param Authorization header string true "token"
// @Param problem_id path string true "题目ID"
// @Success 200 {object} common.GetProblemDetailResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemDetailResponse "1008 需要登录"
// @Failure 200 {object} common.GetProblemDetailResponse "1021 题目ID不存在"
// @Router /problem/{problem_id} [GET]
func (p *ApiProblem) GetProblemDetail(c *gin.Context) {
	var req request.GetProblemDetailReq
	req.ProblemID = c.Param("problem_id")

	data, err := ProblemCache.GetProblemDetailWithCache(req)
	//data, err := ProblemService.GetProblemDetailWithCache(req)
	if err != nil {
		response.ResponseError(c, response.CodeProblemIDNotExist)
		zap.L().Error("controller-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, data)
}

// GetProblemID 获取题目ID接口
// @Tags Problem API
// @Summary 获取题目ID
// @Description 获取题目ID接口
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param title formData string true "标题"
// @Success 200 {object} common.GetProblemIDResponse "1000 获取题目ID成功"
// @Failure 200 {object} common.GetProblemIDResponse "1020 题目title不存在"
// @Failure 200 {object} common.GetProblemIDResponse "1008 需要登录"
// @Router /problem/id [POST]
func (p *ApiProblem) GetProblemID(c *gin.Context) {
	var req request.GetProblemIDReq
	req.Title = c.PostForm("title")

	uid, err := ProblemService.GetProblemID(req)
	if err != nil {
		response.ResponseError(c, response.CodeProblemTitleNotExist)
		return
	}
	response.ResponseSuccess(c, uid)
}

// GetProblemRandom 随机获取一个题目接口
// @Tags Problem API
// @Summary 随机获取一个题目
// @Description 随机获取一个题目接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// // @Param Authorization header string true "token"
// @Success 200 {object} common.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemRandomResponse "1008 需要登录"
// @Router /problem/random [GET]
func (p *ApiProblem) GetProblemRandom(c *gin.Context) {
	data, err := ProblemCache.ProblemIDListCacheRandom()
	if err != nil {
		if err == define.ErrNoProblemIDFound {
			response.ResponseError(c, response.CodeProblemIDNotExist)
			return
		}
		response.ResponseError(c, response.CodeProblemIDNotExist)
		zap.L().Error("controller-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return
	}
	response.ResponseSuccess(c, data)
}

// GetProblemListByCategory 按照分类搜索题目
// @Tags Problem API
// @Summary 按照分类搜索题目
// @Description 按照分类搜索题目
// @Accept json
// @Produce json,multipart/form-data
// @Param category query string true "题目分类"
// @Success 200 {object} common.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemRandomResponse "1001 参数错误"
// @Router /problem/category/search [POST]
func (p *ApiProblem) GetProblemListByCategory(c *gin.Context) {
	category := c.Query("category")

	if category == "" {
		response.ResponseError(c, response.CodeInvalidParam)
		zap.L().Error("controller-GetProblemListByCategory-category is empty ")
		return
	}

	data, err := ProblemService.GetProblemListByCategory(category)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 没有这个分类
			response.ResponseError(c, response.CodeDontHaveThisCategory)
			return
		}
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	if data == nil {
		response.ResponseError(c, response.CodeDontHaveThisCategory)
		return
	}
	response.ResponseSuccess(c, data)
}

// GetCategoryList 获取分类列表
// @Tags Problem API
// @Summary 获取分类列表
// @Description 获取分类列表
// @Accept json
// @Produce json,multipart/form-data
// @Success 200 {object} common.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemRandomResponse "1001 参数错误"
// @Router /problem/category-list [GET]
func (p *ApiProblem) GetCategoryList(c *gin.Context) {
	data, err := ProblemService.GetCategoryList()
	if err != nil {
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	response.ResponseSuccess(c, data)
}

// SearchProblem 搜索题目
// @Tags Problem API
// @Summary 搜索题目
// @Description 搜索题目
// @Accept json
// @Produce json,multipart/form-data
// @Param msg query string true "题目名称"
// @Success 200 {object} common.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemRandomResponse "1001 参数错误"
// @Router /problem/title/search [POST]
func (p *ApiProblem) SearchProblem(c *gin.Context) {
	var req request.SearchProblemReq
	req.Title = c.Query("msg")
	if req.Title == "" {
		response.ResponseError(c, response.CodeInvalidParam)
		zap.L().Error("controller-SearchProblem-msg is empty ")
		return
	}
	data, err := ProblemCache.SearchProblemWithCache(req)
	//data, err := ProblemService.SearchProblem(req)
	if err != nil {
		if err == define.ErrSearchProblem {
			response.ResponseError(c, response.CodeProblemNotFound)
			return
		}
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	response.ResponseSuccess(c, data)
}

// GetHotSearches 获取最热搜索
// @Tags Problem API
// @Summary 获取最热搜索
// @Description 获取最热搜索
// @Accept json
// @Produce json,multipart/form-data
// @Param limit query string false "数量限制"
// @Success 200 {object} common.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemRandomResponse "1014 服务器内部错误"
// @Router /problem/hot-search [GET]
func (p *ApiProblem) GetHotSearches(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", define.DefaultLimit))
	data, err := ProblemCache.GetHotSearches(limit)
	if err != nil {
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	response.ResponseSuccess(c, data)
}

// GetRecentSearches 获取最近搜索
// @Tags Problem API
// @Summary 获取最近搜索
// @Description 获取最近搜索
// @Accept json
// @Produce json,multipart/form-data
// @Param limit query string false "数量限制"
// @Success 200 {object} common.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} common.GetProblemRandomResponse "1014 服务器内部错误"
// @Router /problem/recent-search [GET]
func (p *ApiProblem) GetRecentSearches(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", define.DefaultLimit))
	data, err := ProblemCache.GetRecentSearches(limit)
	if err != nil {
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	response.ResponseSuccess(c, data)
}
