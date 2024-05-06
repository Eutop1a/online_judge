package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	"online-judge/services"
	"strconv"
)

// GetProblemList 获取题目列表接口
// @Tags Problem API
// @Summary 获取题目列表
// @Description 获取题目列表接口
// @Success 200 {object} _Response "获取题目列表成功"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /problem/list [GET]
func GetProblemList(c *gin.Context) {
	var getProblemList services.Problem
	data, err := getProblemList.GetProblemList()
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
// @Param problem_id query string true "题目ID"
// @Success 200 {object} _Response "获取成功"
// @Failure 200 {object} _Response "题目ID不存在"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /problem/{problem_id} [GET]
func GetProblemDetail(c *gin.Context) {
	var getProblemDetail services.Problem
	pid := c.Query("problem_id")
	getProblemDetail.ProblemID = pid

	data, err := getProblemDetail.GetProblemDetail()
	if err != nil {
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
		zap.L().Error("controller-GetProblemDetail-GetProblemDetail ", zap.Error(err))
		return
	}
	resp.ResponseSuccess(c, data)
}

// CreateProblem 创建新题目接口
// @Tags Problem API
// @Summary 创建新题目
// @Description 创建新题目接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param title formData string true "题目标题"
// @Param content formData string true "题目内容"
// @Param difficulty formData string true "题目难度"
// @Param max_runtime formData int true "时间限制"
// @Param max_memory formData int true "内存限制"
// @Param test_cases formData []string true "测试样例集" collectionFormat(multi)
// @Success 200 {object} _Response "创建成功"
// @Failure 200 {object} _Response "参数错误"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /problem/create [POST]
func CreateProblem(c *gin.Context) {
	var createProblem services.Problem

	title := c.PostForm("title")
	content := c.PostForm("content")
	difficulty := c.PostForm("difficulty")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMemory, _ := strconv.Atoi(c.PostForm("max_memory"))

	testCase := c.PostFormArray("test_cases")
	if len(testCase) == 0 {
		zap.L().Error("controller-CreateProblem-PostFormArray testCase is empty")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	//fmt.Println(title)
	//fmt.Println(content)
	//fmt.Println(difficulty)
	//fmt.Println(maxRuntime)
	//fmt.Println(maxMemory)
	//fmt.Println(testCase)
	createProblem.ProblemID = utils.GetUUID()
	createProblem.Content = content
	createProblem.Difficulty = difficulty
	createProblem.Title = title
	createProblem.MaxRuntime = maxRuntime
	createProblem.MaxMemory = maxMemory

	tCase := make([]*services.TestCase, 0)
	for _, value := range testCase {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(value), &caseMap)
		// 检测Map某个键是否存在
		_, iok := caseMap["input"]
		_, ook := caseMap["expected"]
		if err != nil || !iok || !ook {
			resp.ResponseError(c, resp.CodeTestCaseFormatError)
			if err != nil {
				zap.L().Error("controller-CreateProblem-Unmarshal caseMap unmarshal error ", zap.Error(err))
			}
			return
		}
		tCase = append(tCase, &services.TestCase{
			TID:      utils.GetUUID(),
			PID:      createProblem.ProblemID,
			Input:    caseMap["input"],
			Expected: caseMap["expected"],
		})
	}
	createProblem.TestCases = tCase
	response := createProblem.CreateProblem()
	switch response.Code {
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	case resp.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)

	case resp.CreateProblemError:
		resp.ResponseError(c, resp.CodeInternalServerError)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// UpdateProblem 更新题目信息接口
// @Tags Problem API
// @Summary 更新题目信息
// @Description 更新题目信息接口
// @Accept multipart/form-data
// @Produce json
// @Param problem_id query string true "题目ID"
// @Param title formData string false "题目标题"
// @Param content formData string false "题目内容"
// @Param difficulty formData string false "题目难度"
// @Param max_runtime formData string false "时间限制"
// @Param max_memory formData string false "内存限制"
// @Param test_cases formData []string false "测试样例集" collectionFormat(multi)
// @Success 200 {object} _Response "修改成功"
// @Failure 200 {object} _Response "题目ID不存在"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /problem/{problem_id} [PUT]
func UpdateProblem(c *gin.Context) {
	var updateProblem services.Problem

	updateProblem.ProblemID = c.Query("problem_id")

	updateProblem.Title = c.PostForm("title")
	updateProblem.Content = c.PostForm("content")
	updateProblem.Difficulty = c.PostForm("difficulty")
	updateProblem.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	updateProblem.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))
	testCase := c.PostFormArray("test_cases")

	//fmt.Println("id", updateProblem.ProblemID)
	//fmt.Println("title", updateProblem.Title)
	//fmt.Println("content", updateProblem.Content)
	//fmt.Println("difficulty", updateProblem.Difficulty)
	//fmt.Println("max_runtime", updateProblem.MaxRuntime)
	//fmt.Println("max_memory", updateProblem.MaxMemory)
	//fmt.Println("test_cases", testCase)

	tCase := make([]*services.TestCase, 0)
	for _, value := range testCase {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(value), &caseMap)
		// 检测Map某个键是否存在
		_, iok := caseMap["input"]
		_, ook := caseMap["expected"]
		if err != nil || !iok || !ook {
			resp.ResponseError(c, resp.CodeTestCaseFormatError)
			if err != nil {
				zap.L().Error("controller-UpdateProblem-Unmarshal caseMap unmarshal error ", zap.Error(err))
			}
			return
		}
		tCase = append(tCase, &services.TestCase{
			TID:      utils.GetUUID(),
			PID:      updateProblem.ProblemID,
			Input:    caseMap["input"],
			Expected: caseMap["expected"],
		})
	}
	updateProblem.TestCases = tCase
	response := updateProblem.UpdateProblem()
	switch response.Code {
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)
	case resp.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
	case resp.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}

// DeleteProblem 删除题目接口
// @Tags Problem API
// @Summary 删除题目
// @Description 删除题目接口
// @Accept multipart/form-data
// @Produce json
// @Param problem_id query string true "题目ID"
// @Success 200 {object} _Response "删除成功"
// @Failure 200 {object} _Response "题目ID不存在"
// @Failure 200 {object} _Response "服务器内部错误"
// @Router /problem/{problem_id} [DELETE]
func DeleteProblem(c *gin.Context) {
	var deleteProblem services.Problem
	id := c.Query("problem_id")

	deleteProblem.ProblemID = id

	response := deleteProblem.DeleteProblem()
	switch response.Code {
	case resp.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	case resp.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}

// GetProblemID 获取题目ID接口
// @Tags Problem API
// @Summary 获取题目ID
// @Description 获取题目ID接口
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "标题"
// @Success 200 {object} _Response "获取题目ID成功"
// @Failure 200 {object} _Response "题目title不存在"
// @Router /problem/id [POST]
func GetProblemID(c *gin.Context) {
	var getProblemID services.Problem
	title := c.PostForm("title")
	getProblemID.Title = title
	uid, err := getProblemID.GetProblemID()
	if err != nil {
		resp.ResponseError(c, resp.CodeProblemTitleNotExist)
		zap.L().Error("controller-GetProblemID-GetProblemID ", zap.Error(err))
		return
	}
	resp.ResponseSuccess(c, uid)
}
