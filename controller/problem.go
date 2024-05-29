package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online-judge/dao/redis"
	"online-judge/pkg/define"
	"online-judge/pkg/resp"
	"online-judge/services"
	"strconv"
)

// GetProblemList 获取题目列表接口
// @Tags Problem API
// @Summary 获取题目列表
// @Description 获取题目列表接口
// @Param Authorization header string true "token"
// @Param page query int false "input current page num, default: 1"
// @Param size query int false "pageSize"
// @Success 200 {object} models.GetProblemListResponse "获取题目列表成功"
// @Failure 200 {object} models.GetProblemListResponse "需要登录"
// @Failure 200 {object} models.GetProblemListResponse "服务器内部错误"
// @Router /problem/list [GET]
func GetProblemList(c *gin.Context) {
	var getProblemList services.Problem
	getProblemList.Size, _ = strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	getProblemList.Page, _ = strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	data, err := getProblemList.GetProblemListWithCache(redis.Client)
	//data, err := getProblemList.GetProblemList()
	if err != nil {
		resp.ResponseError(c, resp.CodeInternalServerError)
		zap.L().Error("controller-GetProblemList-GetProblemList ", zap.Error(err))
		return
	}
	resp.ResponseSuccess(c, data)
}

// GetProblemDetail 获取单个题目详细接口
// @Tags Problem API
// @Summary 获取单个题目详细
// @Description 获取单个题目详细接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param problem_id path string true "题目ID"
// @Success 200 {object} models.GetProblemDetailResponse "1000 获取成功"
// @Failure 200 {object} models.GetProblemDetailResponse "1008 需要登录"
// @Failure 200 {object} models.GetProblemDetailResponse "1021 题目ID不存在"
// @Router /problem/{problem_id} [GET]
func GetProblemDetail(c *gin.Context) {
	var getProblemDetail services.Problem
	pid := c.Param("problem_id")
	getProblemDetail.ProblemID = pid

	data, err := getProblemDetail.GetProblemDetailWithCache(redis.Client)
	if err != nil {
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
		zap.L().Error("controller-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return
	}
	resp.ResponseSuccess(c, data)
}

// GetProblemID 获取题目ID接口
// @Tags Problem API
// @Summary 获取题目ID
// @Description 获取题目ID接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param title formData string true "标题"
// @Success 200 {object} models.GetProblemIDResponse "1000 获取题目ID成功"
// @Failure 200 {object} models.GetProblemIDResponse "1020 题目title不存在"
// @Failure 200 {object} models.GetProblemIDResponse "1008 需要登录"
// @Router /problem/id [POST]
func GetProblemID(c *gin.Context) {
	var getProblemID services.Problem
	title := c.PostForm("title")
	getProblemID.Title = title
	uid, err := getProblemID.GetProblemID()
	if err != nil {
		resp.ResponseError(c, resp.CodeProblemTitleNotExist)
		return
	}
	resp.ResponseSuccess(c, uid)
}

// GetProblemRandom 随机获取一个题目接口
// @Tags Problem API
// @Summary 随机获取一个题目
// @Description 随机获取一个题目接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Success 200 {object} models.GetProblemRandomResponse "1000 获取成功"
// @Failure 200 {object} models.GetProblemRandomResponse "1008 需要登录"
// @Router /problem/random [GET]
func GetProblemRandom(c *gin.Context) {
	var getProblemDetail services.Problem

	data, err := getProblemDetail.GetProblemRandom(redis.Client)
	if err != nil {
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
		zap.L().Error("controller-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return
	}
	resp.ResponseSuccess(c, data)
}
