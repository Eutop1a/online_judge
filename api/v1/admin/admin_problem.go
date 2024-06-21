package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"online_judge/consts"
	"online_judge/consts/resp_code"
	"online_judge/dao/redis"
	"online_judge/models/admin/request"
	"online_judge/models/common/response"
	"online_judge/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
)

type ApiAdminProblem struct{}

/*
// CreateProblem 创建新题目接口
// @Tags Admin API
// @Summary 创建新题目
// @Description 创建新题目接口
// @Accept application/json,multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param title formData string true "题目标题"
// @Param category formData []string true "分类id"  collectionFormat(multi)
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
*/

// CreateProblem 创建新题目接口
// @swagger:order
// @Tags Admin API
// @Summary 创建新题目
// @Description 创建新题目接口
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param req body request.AdminCreateProblemReq true "创建题目信息的请求体"
// @Success 200 {object} common.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} common.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} common.CreateProblemResponse "1018 测试用例格式错误"
// @Failure 200 {object} common.CreateProblemResponse "1019 题目标题已存在"
// @Failure 200 {object} common.CreateProblemResponse "1008 需要登录"
// @Failure 200 {object} common.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/problem/create [POST]
func (a *ApiAdminProblem) CreateProblem(c *gin.Context) {
	//var createProblem services.Problem
	// 解析 JSON 请求体
	var req request.AdminCreateProblemReq

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("controller-CreateProblem-BindJSON error", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	// 测试用例为空
	if len(req.TestCases) == 0 {
		zap.L().Error("controller-CreateProblem-TestCases is empty")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	// 题目类型为空
	if len(req.Category) == 0 {
		zap.L().Error("controller-CreateProblem-Category is empty")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	req.ProblemID = utils.GetUUID()

	for _, v := range req.TestCases {
		v.TID = utils.GetUUID()
		v.PID = req.ProblemID
	}
	req.RedisClient = redis.Client
	req.Ctx = redis.Ctx

	resp := AdminService.CreateProblem(req)
	switch resp.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.ProblemAlreadyExist:
		response.ResponseError(c, response.CodeProblemTitleExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

/*
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
// @Param category formData []string false "分类id" collectionFormat(multi)
// @Param test_cases formData []string false "测试样例集" collectionFormat(multi)
// @Success 200 {object} common.UpdateProblemResponse "修改成功"
// @Failure 200 {object} common.UpdateProblemResponse "题目ID不存在"
// @Failure 200 {object} common.UpdateProblemResponse "题目标题已存在"
// @Failure 200 {object} common.UpdateProblemResponse "测试用例格式错误"
// @Failure 200 {object} common.UpdateProblemResponse "需要登录"
// @Failure 200 {object} common.UpdateProblemResponse "服务器内部错误"
// @Router /admin/problem/{problem_id} [PUT]
*/

// UpdateProblem 更新题目信息接口
// @Tags Admin API
// @Summary 更新题目信息
// @Description 更新题目信息接口
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param problem_id query string true "题目ID"
// @Param req body request.AdminUpdateProblemReq true "更新题目信息的请求体"
// @Success 200 {object} common.UpdateProblemResponse "修改成功"
// @Failure 200 {object} common.UpdateProblemResponse "题目ID不存在"
// @Failure 200 {object} common.UpdateProblemResponse "题目标题已存在"
// @Failure 200 {object} common.UpdateProblemResponse "测试用例格式错误"
// @Failure 200 {object} common.UpdateProblemResponse "需要登录"
// @Failure 200 {object} common.UpdateProblemResponse "服务器内部错误"
// @Router /admin/problem/update [PUT]
func (a *ApiAdminProblem) UpdateProblem(c *gin.Context) {

	var req request.AdminUpdateProblemReq
	req.ProblemID = c.Query("problem_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("controller-UpdateProblem-BindJSON error", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	if len(req.TestCases) == 0 {
		zap.L().Error("controller-UpdateProblem-TestCases is empty")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	for _, v := range req.TestCases {
		v.TID = utils.GetUUID()
		v.PID = req.ProblemID
	}

	req.RedisClient = redis.Client
	req.Ctx = redis.Ctx

	resp := AdminService.UpdateProblem(req)
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
// @Param problem_id query string true "题目ID"
// @Success 200 {object} common.DeleteProblemResponse "删除成功"
// @Failure 200 {object} common.DeleteProblemResponse "题目ID不存在"
// @Failure 200 {object} common.DeleteProblemResponse "需要登录"
// @Failure 200 {object} common.DeleteProblemResponse "服务器内部错误"
// @Router /admin/problem/delete [DELETE]
func (a *ApiAdminProblem) DeleteProblem(c *gin.Context) {
	var req request.AdminDeleteProblemReq
	req.ProblemID = c.Query("problem_id")
	req.RedisClient = redis.Client
	req.Ctx = redis.Ctx

	resp := AdminService.DeleteProblem(req)

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
func (a *ApiAdminProblem) CreateProblemWithFile(c *gin.Context) {
	var req request.AdminCreateProblemWithFileReq

	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		zap.L().Error("controller-CreateProblemWithFile-ParseMultipartForm ", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	req.Title = c.PostForm("title")
	req.Content = c.PostForm("content")
	req.Difficulty = c.PostForm("difficulty")
	req.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	req.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	inputFile := c.Request.MultipartForm.File["input"]
	outputFile := c.Request.MultipartForm.File["expected"]
	req.ProblemID = utils.GetUUID()

	problemDir := filepath.Join(consts.FilePath, req.ProblemID)
	err = os.MkdirAll(problemDir, os.ModePerm)
	if err != nil {
		zap.L().Error("controller-CreateProblemWithFile-MkdirAll", zap.Error(err))
		response.ResponseError(c, response.CodeInternalServerError)
		return
	}

	inputDst := filepath.Join(consts.FilePath, req.ProblemID, consts.FileInput)
	expectedDst := filepath.Join(consts.FilePath, req.ProblemID, consts.FileExpected)

	req.InputDst = inputDst
	req.ExpectedDst = expectedDst

	// 保存输入文件
	a.SaveFile(c, inputFile, inputDst)
	// 保存输出文件
	a.SaveFile(c, outputFile, expectedDst)

	resp := AdminService.CreateProblemWithFile(req)
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
func (a *ApiAdminProblem) DeleteProblemWithFile(c *gin.Context) {
	var req request.AdminDeleteProblemWithFileReq
	req.ProblemID = c.Param("problem_id")

	resp := AdminService.DeleteProblemWithFile(req)
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
func (a *ApiAdminProblem) UpdateProblemWithFile(c *gin.Context) {
	var req request.AdminUpdateProblemWithFileReq

	req.ProblemID = c.PostForm("problem_id")
	req.Title = c.PostForm("title")
	req.Content = c.PostForm("content")
	req.Difficulty = c.PostForm("difficulty")
	req.MaxRuntime, _ = strconv.Atoi(c.PostForm("max_runtime"))
	req.MaxMemory, _ = strconv.Atoi(c.PostForm("max_memory"))

	inputFile := c.Request.MultipartForm.File["input"]
	outputFile := c.Request.MultipartForm.File["expected"]

	req.InputDst = filepath.Join(consts.FilePath, req.ProblemID, consts.FileInput)
	req.ExpectedDst = filepath.Join(consts.FilePath, req.ProblemID, consts.FileExpected)
	// 在添加之前先删除所有的测试用例
	resp := AdminService.DeleteProblemTestCaseWithFile(req)
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
	a.SaveFile(c, inputFile, req.InputDst)
	// 保存输出文件
	a.SaveFile(c, outputFile, req.ExpectedDst)

	resp = AdminService.UpdateProblemWithFile(req)
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

func (a *ApiAdminProblem) SaveFile(c *gin.Context, fileHeader []*multipart.FileHeader, dstDir string) {
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
