package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"online_judge/consts"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/dao/redis"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"online_judge/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
)

// CreateProblem 创建新题目接口
// @Tags Admin API
// @Summary 创建新题目
// @Description 创建新题目接口
// @Accept application/json,multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param title formData string true "题目标题"
// @Param content formData string true "题目内容"
// @Param difficulty formData string true "题目难度"
// @Param max_runtime formData int true "时间限制"
// @Param max_memory formData int true "内存限制"
// @Param test_cases formData []string true "测试样例集" collectionFormat(multi)
// @Success 200 {object} common.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} common.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} common.CreateProblemResponse "1018 测试用例格式错误"
// @Failure 200 {object} common.CreateProblemResponse "1019 题目标题已存在"
// @Failure 200 {object} common.CreateProblemResponse "1008 需要登录"
// @Failure 200 {object} common.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/problem/create [POST]
func (a *ApiAdmin) CreateProblem(c *gin.Context) {
	//var createProblem services.Problem
	// 解析 JSON 请求体
	var createProblemReq request.AdminCreateProblemReq
	if err := c.ShouldBind(&createProblemReq); err != nil {
		zap.L().Error("controller-CreateProblem-BindJSON error", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	if len(createProblemReq.TestCases) == 0 {
		zap.L().Error("controller-CreateProblem-TestCases is empty")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	createProblem := mysql.Problems{
		ProblemID:  utils.GetUUID(),
		Title:      createProblemReq.Title,
		Content:    createProblemReq.Content,
		Difficulty: createProblemReq.Difficulty,
		MaxRuntime: createProblemReq.MaxRuntime,
		MaxMemory:  createProblemReq.MaxMemory,
		TestCases:  make([]*mysql.TestCase, len(createProblemReq.TestCases)),
	}

	for i, tc := range createProblemReq.TestCases {
		createProblem.TestCases[i] = &mysql.TestCase{
			TID:      utils.GetUUID(),
			PID:      createProblem.ProblemID,
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	createProblemReq.RedisClient = redis.Client
	createProblemReq.Ctx = redis.Ctx

	resp := AdminService.CreateProblem(createProblemReq)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.ProblemAlreadyExist:
		response.ResponseError(c, response.CodeProblemTitleExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// UpdateProblem 更新题目信息接口
// @Tags Admin API
// @Summary 更新题目信息
// @Description 更新题目信息接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id path string true "题目ID"
// @Param title formData string false "题目标题"
// @Param content formData string false "题目内容"
// @Param difficulty formData string false "题目难度"
// @Param max_runtime formData string false "时间限制"
// @Param max_memory formData string false "内存限制"
// @Param test_cases formData []string false "测试样例集" collectionFormat(multi)
// @Success 200 {object} common.UpdateProblemResponse "修改成功"
// @Failure 200 {object} common.UpdateProblemResponse "题目ID不存在"
// @Failure 200 {object} common.UpdateProblemResponse "题目标题已存在"
// @Failure 200 {object} common.UpdateProblemResponse "测试用例格式错误"
// @Failure 200 {object} common.UpdateProblemResponse "需要登录"
// @Failure 200 {object} common.UpdateProblemResponse "服务器内部错误"
// @Router /admin/problem/{problem_id} [PUT]
func (a *ApiAdmin) UpdateProblem(c *gin.Context) {

	var updateProblemReq request.AdminUpdateProblemReq
	if err := c.ShouldBindJSON(&updateProblemReq); err != nil {
		zap.L().Error("controller-UpdateProblem-BindJSON error", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	if len(updateProblemReq.TestCases) == 0 {
		zap.L().Error("controller-UpdateProblem-TestCases is empty")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	updateProblem := mysql.Problems{
		ProblemID:  c.Param("problem_id"),
		Title:      updateProblemReq.Title,
		Content:    updateProblemReq.Content,
		Difficulty: updateProblemReq.Difficulty,
		MaxRuntime: updateProblemReq.MaxRuntime,
		MaxMemory:  updateProblemReq.MaxMemory,
		TestCases:  make([]*mysql.TestCase, len(updateProblemReq.TestCases)),
	}

	for i, tc := range updateProblemReq.TestCases {
		updateProblem.TestCases[i] = &mysql.TestCase{
			TID:      utils.GetUUID(),
			PID:      updateProblem.ProblemID,
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}

	updateProblemReq.RedisClient = redis.Client
	updateProblemReq.Ctx = redis.Ctx

	resp := AdminService.UpdateProblem(updateProblemReq)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)
	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)
	case resp_code.ProblemAlreadyExist:
		response.ResponseError(c, response.CodeProblemTitleExist)
	default:
		response.ResponseError(c, response.CodeInternalServerError)
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
// @Success 200 {object} common.DeleteProblemResponse "删除成功"
// @Failure 200 {object} common.DeleteProblemResponse "题目ID不存在"
// @Failure 200 {object} common.DeleteProblemResponse "需要登录"
// @Failure 200 {object} common.DeleteProblemResponse "服务器内部错误"
// @Router /admin/problem/{problem_id} [DELETE]
func (a *ApiAdmin) DeleteProblem(c *gin.Context) {
	var deleteProblemReq request.AdminDeleteProblemReq
	deleteProblemReq.ProblemID = c.Param("problem_id")

	deleteProblemReq.RedisClient = redis.Client
	deleteProblemReq.Ctx = redis.Ctx

	resp := AdminService.DeleteProblem(deleteProblemReq)

	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
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
// @Success 200 {object} common.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} common.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} common.CreateProblemResponse "1018 测试用例格式错误"
// @Failure 200 {object} common.CreateProblemResponse "1019 题目标题已存在"
// @Failure 200 {object} common.CreateProblemResponse "1008 需要登录"
// @Failure 200 {object} common.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/problem/file/create [POST]
func (a *ApiAdmin) CreateProblemWithFile(c *gin.Context) {
	var createProblemReq request.AdminCreateProblemWithFileReq

	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		zap.L().Error("controller-CreateProblemWithFile-ParseMultipartForm ", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	createProblemReq.Title = c.PostForm("title")
	createProblemReq.Content = c.PostForm("content")
	createProblemReq.Difficulty = c.PostForm("difficulty")
	createProblemReq.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	createProblemReq.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	inputFile := c.Request.MultipartForm.File["input"]
	outputFile := c.Request.MultipartForm.File["expected"]
	createProblemReq.ProblemID = utils.GetUUID()

	problemDir := filepath.Join(consts.FilePath, createProblemReq.ProblemID)
	err = os.MkdirAll(problemDir, os.ModePerm)
	if err != nil {
		zap.L().Error("controller-CreateProblemWithFile-MkdirAll", zap.Error(err))
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}

	inputDst := filepath.Join(consts.FilePath, createProblemReq.ProblemID, consts.FileInput)
	expectedDst := filepath.Join(consts.FilePath, createProblemReq.ProblemID, consts.FileExpected)

	createProblemReq.InputDst = inputDst
	createProblemReq.ExpectedDst = expectedDst

	// 保存输入文件
	a.SaveFile(c, inputFile, inputDst)
	// 保存输出文件
	a.SaveFile(c, outputFile, expectedDst)

	resp := AdminService.CreateProblemWithFile(createProblemReq)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.ProblemAlreadyExist:
		response.ResponseError(c, response.CodeProblemTitleExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
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
// @Success 200 {object} common.DeleteProblemResponse "删除成功"
// @Failure 200 {object} common.DeleteProblemResponse "题目ID不存在"
// @Failure 200 {object} common.DeleteProblemResponse "需要登录"
// @Failure 200 {object} common.DeleteProblemResponse "服务器内部错误"
// @Router /admin/problem/file/{problem_id} [DELETE]
func (a *ApiAdmin) DeleteProblemWithFile(c *gin.Context) {
	var deleteProblemReq request.AdminDeleteProblemWithFileReq
	deleteProblemReq.ProblemID = c.Param("problem_id")

	resp := AdminService.DeleteProblemWithFile(deleteProblemReq)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
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
// @Success 200 {object} common.UpdateProblemResponse "修改成功"
// @Failure 200 {object} common.UpdateProblemResponse "题目ID不存在"
// @Failure 200 {object} common.UpdateProblemResponse "题目标题已存在"
// @Failure 200 {object} common.UpdateProblemResponse "测试用例格式错误"
// @Failure 200 {object} common.UpdateProblemResponse "需要登录"
// @Failure 200 {object} common.UpdateProblemResponse "服务器内部错误"
// @Router /admin/problem/file/update [PUT]
func (a *ApiAdmin) UpdateProblemWithFile(c *gin.Context) {
	var updateProblemReq request.AdminUpdateProblemWithFileReq

	updateProblemReq.ProblemID = c.PostForm("problem_id")
	updateProblemReq.Title = c.PostForm("title")
	updateProblemReq.Content = c.PostForm("content")
	updateProblemReq.Difficulty = c.PostForm("difficulty")
	updateProblemReq.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	updateProblemReq.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	inputFile := c.Request.MultipartForm.File["input"]
	outputFile := c.Request.MultipartForm.File["expected"]

	updateProblemReq.InputDst = filepath.Join(consts.FilePath, updateProblemReq.ProblemID, consts.FileInput)
	updateProblemReq.ExpectedDst = filepath.Join(consts.FilePath, updateProblemReq.ProblemID, consts.FileExpected)
	// 在添加之前先删除所有的测试用例
	resp := AdminService.DeleteProblemTestCaseWithFile(updateProblemReq)
	switch resp.Code {
	case resp_code.Success:

	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)
		return

	case resp_code.ProblemAlreadyExist:
		response.ResponseError(c, response.CodeProblemTitleExist)
		return

	default:
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	// 保存输入文件
	a.SaveFile(c, inputFile, updateProblemReq.InputDst)
	// 保存输出文件
	a.SaveFile(c, outputFile, updateProblemReq.ExpectedDst)

	resp = AdminService.UpdateProblemWithFile(updateProblemReq)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)
	case resp_code.ProblemNotExist:
		response.ResponseError(c, response.CodeProblemIDNotExist)
	case resp_code.ProblemAlreadyExist:
		response.ResponseError(c, response.CodeProblemTitleExist)
	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
	return
}

func (a *ApiAdmin) SaveFile(c *gin.Context, fileHeader []*multipart.FileHeader, dstDir string) {
	// 保存输出文件
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		zap.L().Error("controller-SaveFile-MkdirAll", zap.Error(err))
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}
	for _, file := range fileHeader {
		dst := filepath.Join(dstDir, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			zap.L().Error("controller-SaveFile-SaveUploadedFile", zap.Error(err))
			response.ResponseError(c, response.CodeInternalServerError)
			return
		}
	}
}
