package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online-judge/consts"
	"online-judge/consts/resp_code"
	"online-judge/dao/redis"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	"online-judge/services"
	"os"
	"path/filepath"
	"strconv"
)

// AddSuperAdmin 添加超级管理员接口
// @Tags Admin API
// @Summary 添加超级管理员
// @Description 添加超级管理员接口
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param secret formData string true "密钥"
// @Success 200 {object} models.AddSuperAdminResponse "添加超级管理员成功"
// @Failure 200 {object} models.AddSuperAdminResponse "参数错误"
// @Failure 200 {object} models.AddSuperAdminResponse "没有此用户ID"
// @Failure 200 {object} models.AddSuperAdminResponse "用户已是管理员"
// @Failure 200 {object} models.AddSuperAdminResponse "密钥错误"
// @Failure 200 {object} models.AddSuperAdminResponse "服务器内部错误"
// @Router /admin/users/add-super-admin [POST]
func AddSuperAdmin(c *gin.Context) {
	var addAdmin services.UserService
	//uid := c.PostForm("user_id")
	username := c.PostForm("username")
	secret := c.PostForm("secret")
	if username == "" {
		zap.L().Error("controller-AddSuperAdmin-PostForm add admin params error")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	//addAdmin.UserID, _ = strconv.ParseInt(uid, 10, 64)
	addAdmin.UserName = username
	var ret resp.Response
	ret = addAdmin.AddSuperAdmin(secret)
	switch ret.Code {

	case resp_code.Success:
		resp.ResponseSuccess(c, resp_code.Success)

	case resp_code.NotExistUsername:
		resp.ResponseError(c, resp.CodeUsernameNotExist)

	case resp_code.SecretError:
		resp.ResponseError(c, resp.CodeErrorSecret)

	case resp_code.UsernameAlreadyExist:
		resp.ResponseError(c, resp.CodeUsernameAlreadyExist)

	default:
		resp.ResponseError(c, resp.CodeInvalidParam)
	}
	return
}

// DeleteUser 删除用户接口
// @Tags Admin API
// @Summary 删除用户
// @Description 删除用户接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param user_id path string true "用户ID"
// @Success 200 {object} models.DeleteUserResponse "删除用户成功"
// @Failure 200 {object} models.DeleteUserResponse "参数错误"
// @Failure 200 {object} models.DeleteUserResponse "没有此用户ID"
// @Failure 200 {object} models.DeleteUserResponse "需要登录"
// @Failure 200 {object} models.DeleteUserResponse "服务器内部错误"
// @Router /admin/users/{user_id} [DELETE]
func DeleteUser(c *gin.Context) {
	var deleteUser services.UserService
	uid := c.Param("user_id")
	if uid == "" {
		zap.L().Error("controller-DeleteUser-Param deleteUser params error")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	deleteUser.UserID, _ = strconv.ParseInt(uid, 10, 64)
	var ret resp.Response
	ret = deleteUser.DeleteUser()

	switch ret.Code {
	// 成功
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 用户不存在
	case resp_code.NotExistUserID:
		resp.ResponseError(c, resp.CodeUseNotExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// AddAdmin 添加管理员接口
// @Tags Admin API
// @Summary 添加管理员
// @Description 添加管理员接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param username formData string true "用户名"
// @Success 200 {object} models.AddAdminResponse "删除用户成功"
// @Failure 200 {object} models.AddAdminResponse "参数错误"
// @Failure 200 {object} models.AddAdminResponse "没有此用户ID"
// @Failure 200 {object} models.AddAdminResponse "需要登录"
// @Failure 200 {object} models.AddAdminResponse "服务器内部错误"
// @Router /admin/users/add-admin [POST]
func AddAdmin(c *gin.Context) {
	var addAdmin services.UserService
	//uid := c.PostForm("user_id")
	username := c.PostForm("username")
	if username == "" {
		zap.L().Error("add admin params error")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}
	//addAdmin.UserID, _ = strconv.ParseInt(uid, 10, 64)
	addAdmin.UserName = username
	var ret resp.Response
	ret = addAdmin.AddAdmin()

	switch ret.Code {
	// 成功
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	// 用户名不存在
	case resp_code.NotExistUsername:
		resp.ResponseError(c, resp.CodeUsernameNotExist)

	// 用户已经是管理员了
	case resp_code.UsernameAlreadyExist:
		resp.ResponseError(c, resp.CodeUsernameAlreadyExist)

	// 服务器内部错误
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// CreateProblem 创建新题目接口
// @Tags Admin API
// @Summary 创建新题目
// @Description 创建新题目接口
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param title formData string true "题目标题"
// @Param content formData string true "题目内容"
// @Param difficulty formData string true "题目难度"
// @Param max_runtime formData int true "时间限制"
// @Param max_memory formData int true "内存限制"
// @Param test_cases formData []string true "测试样例集" collectionFormat(multi)
// @Success 200 {object} models.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} models.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} models.CreateProblemResponse "1018 测试用例格式错误"
// @Failure 200 {object} models.CreateProblemResponse "1019 题目标题已存在"
// @Failure 200 {object} models.CreateProblemResponse "1008 需要登录"
// @Failure 200 {object} models.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/problem/create [POST]
func CreateProblem(c *gin.Context) {
	var createProblem services.Problem

	createProblem.Title = c.PostForm("title")
	createProblem.Content = c.PostForm("content")
	createProblem.Difficulty = c.PostForm("difficulty")
	createProblem.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	createProblem.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	testCase := c.PostFormArray("test_cases")
	if len(testCase) == 0 {
		zap.L().Error("controller-CreateProblem-PostFormArray testCase is empty")
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}

	createProblem.ProblemID = utils.GetUUID()

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

	response := createProblem.CreateProblem(redis.Client, redis.Ctx)
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	case resp_code.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// UpdateProblem 更新题目信息接口
// @Tags Admin API
// @Summary 更新题目信息
// @Description 更新题目信息接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id formData string true "题目ID"
// @Param title formData string false "题目标题"
// @Param content formData string false "题目内容"
// @Param difficulty formData string false "题目难度"
// @Param max_runtime formData string false "时间限制"
// @Param max_memory formData string false "内存限制"
// @Param test_cases formData []string false "测试样例集" collectionFormat(multi)
// @Success 200 {object} models.UpdateProblemResponse "修改成功"
// @Failure 200 {object} models.UpdateProblemResponse "题目ID不存在"
// @Failure 200 {object} models.UpdateProblemResponse "题目标题已存在"
// @Failure 200 {object} models.UpdateProblemResponse "测试用例格式错误"
// @Failure 200 {object} models.UpdateProblemResponse "需要登录"
// @Failure 200 {object} models.UpdateProblemResponse "服务器内部错误"
// @Router /admin/problem/{problem_id} [PUT]
func UpdateProblem(c *gin.Context) {
	var updateProblem services.Problem

	updateProblem.ProblemID = c.PostForm("problem_id")
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
	response := updateProblem.UpdateProblem(redis.Client, redis.Ctx)
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)
	case resp_code.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
	case resp_code.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}

// DeleteProblem 删除题目接口
// @Tags Admin API
// @Summary 删除题目
// @Description 删除题目接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id path string true "题目ID"
// @Success 200 {object} models.DeleteProblemResponse "删除成功"
// @Failure 200 {object} models.DeleteProblemResponse "题目ID不存在"
// @Failure 200 {object} models.DeleteProblemResponse "需要登录"
// @Failure 200 {object} models.DeleteProblemResponse "服务器内部错误"
// @Router /admin/problem/{problem_id} [DELETE]
func DeleteProblem(c *gin.Context) {
	var deleteProblem services.Problem
	deleteProblem.ProblemID = c.Param("problem_id")
	response := deleteProblem.DeleteProblem(redis.Client, redis.Ctx)
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	case resp_code.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}

// CreateProblemWithFile 创建新题目接口，输入输出是文件的形式
// @Tags Admin API
// @Summary 创建新题目，输入输出是文件的形式
// @Description 创建新题目接口，输入输出是文件的形式
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param title formData string true "题目标题"
// @Param content formData string true "题目内容"
// @Param difficulty formData string true "题目难度"
// @Param max_runtime formData int true "时间限制"
// @Param max_memory formData int true "内存限制"
// @Param input formData []file true "问题的输入文件(.in)" collectionFormat(multi)
// @Param expected formData []file true "问题的输出文件(.out)" collectionFormat(multi)
// @Success 200 {object} models.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} models.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} models.CreateProblemResponse "1018 测试用例格式错误"
// @Failure 200 {object} models.CreateProblemResponse "1019 题目标题已存在"
// @Failure 200 {object} models.CreateProblemResponse "1008 需要登录"
// @Failure 200 {object} models.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/problem/file/create [POST]
func CreateProblemWithFile(c *gin.Context) {
	var createProblem services.Problem

	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		zap.L().Error("controller-CreateProblemWithFile-ParseMultipartForm ", zap.Error(err))
		resp.ResponseError(c, resp.CodeInvalidParam)
		return
	}

	createProblem.Title = c.PostForm("title")
	createProblem.Content = c.PostForm("content")
	createProblem.Difficulty = c.PostForm("difficulty")
	createProblem.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	createProblem.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	inputFile := c.Request.MultipartForm.File["input"]
	outputFile := c.Request.MultipartForm.File["expected"]
	createProblem.ProblemID = utils.GetUUID()

	problemDir := filepath.Join(consts.FilePath, createProblem.ProblemID)
	err = os.MkdirAll(problemDir, os.ModePerm)
	if err != nil {
		zap.L().Error("controller-CreateProblemWithFile-MkdirAll", zap.Error(err))
		resp.ResponseError(c, resp.CodeInternalServerError)
		return
	}

	inputDst := filepath.Join(consts.FilePath, createProblem.ProblemID, consts.FileInput)
	expectedDst := filepath.Join(consts.FilePath, createProblem.ProblemID, consts.FileExpected)

	createProblem.InputDst = inputDst
	createProblem.ExpectedDst = expectedDst

	// 保存输入文件
	SaveFile(c, inputFile, inputDst)
	// 保存输出文件
	SaveFile(c, outputFile, expectedDst)

	response := createProblem.CreateProblemWithFile()
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	case resp_code.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
}

// DeleteProblemWithFile 删除题目接口，输入输出是文件的形式
// @Tags Admin API
// @Summary 删除题目，输入输出是文件的形式
// @Description 删除题目接口，输入输出是文件的形式
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id path string true "题目ID"
// @Success 200 {object} models.DeleteProblemResponse "删除成功"
// @Failure 200 {object} models.DeleteProblemResponse "题目ID不存在"
// @Failure 200 {object} models.DeleteProblemResponse "需要登录"
// @Failure 200 {object} models.DeleteProblemResponse "服务器内部错误"
// @Router /admin/problem/file/{problem_id} [DELETE]
func DeleteProblemWithFile(c *gin.Context) {
	var deleteProblem services.Problem
	deleteProblem.ProblemID = c.Param("problem_id")

	response := deleteProblem.DeleteProblemWithFile()
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)

	case resp_code.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}

// UpdateProblemWithFile 更新题目信息接口
// @Tags Admin API
// @Summary 更新题目信息
// @Description 更新题目信息接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id formData string true "题目ID"
// @Param title formData string false "题目标题"
// @Param content formData string false "题目内容"
// @Param difficulty formData string false "题目难度"
// @Param max_runtime formData string false "时间限制"
// @Param max_memory formData string false "内存限制"
// @Param input formData []file true "问题的输入文件(.in)" collectionFormat(multi)
// @Param expected formData []file true "问题的输出文件(.out)" collectionFormat(multi)
// @Success 200 {object} models.UpdateProblemResponse "修改成功"
// @Failure 200 {object} models.UpdateProblemResponse "题目ID不存在"
// @Failure 200 {object} models.UpdateProblemResponse "题目标题已存在"
// @Failure 200 {object} models.UpdateProblemResponse "测试用例格式错误"
// @Failure 200 {object} models.UpdateProblemResponse "需要登录"
// @Failure 200 {object} models.UpdateProblemResponse "服务器内部错误"
// @Router /admin/problem/file/update [PUT]
func UpdateProblemWithFile(c *gin.Context) {
	var updateProblem services.Problem

	updateProblem.ProblemID = c.PostForm("problem_id")
	updateProblem.Title = c.PostForm("title")
	updateProblem.Content = c.PostForm("content")
	updateProblem.Difficulty = c.PostForm("difficulty")
	updateProblem.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	updateProblem.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	inputFile := c.Request.MultipartForm.File["input"]
	outputFile := c.Request.MultipartForm.File["expected"]

	updateProblem.InputDst = filepath.Join(consts.FilePath, updateProblem.ProblemID, consts.FileInput)
	updateProblem.ExpectedDst = filepath.Join(consts.FilePath, updateProblem.ProblemID, consts.FileExpected)
	// 在添加之前先删除所有的测试用例
	response := updateProblem.DeleteProblemTestCaseWithFile()
	switch response.Code {
	case resp_code.Success:

	case resp_code.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
		return

	case resp_code.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)
		return

	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
		return
	}
	// 保存输入文件
	SaveFile(c, inputFile, updateProblem.InputDst)
	// 保存输出文件
	SaveFile(c, outputFile, updateProblem.ExpectedDst)

	response = updateProblem.UpdateProblemWithFile()
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, resp.CodeSuccess)
	case resp_code.ProblemNotExist:
		resp.ResponseError(c, resp.CodeProblemIDNotExist)
	case resp_code.ProblemAlreadyExist:
		resp.ResponseError(c, resp.CodeProblemTitleExist)
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}
