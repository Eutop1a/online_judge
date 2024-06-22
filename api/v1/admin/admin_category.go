package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/models/common/response"
)

type ApiAdminCategory struct{}

// AddCategory 增加分类
// @Tags Admin API
// @Summary 增加分类
// @Description 增加分类接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param category_name formData string true "分类名称"
// @Success 200 {object} common.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} common.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} common.CreateProblemResponse "1033 题目分类已存在"
// @Failure 200 {object} common.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/category/create [POST]
func (a *ApiAdminCategory) AddCategory(c *gin.Context) {
	// 分类名称
	name := c.PostForm("category_name")
	if name == "" {
		zap.L().Error("category name should not be empty")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	ret := AdminService.AddCategory(name)
	switch ret.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.CategoryTypeAlreadyExist:
		response.ResponseError(c, response.CodeCategoryTypeAlreadyExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// UpdateCategory 更新分类
// @Tags Admin API
// @Summary 更新分类
// @Description 更新分类接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param category_id formData string true "分类ID"
// @Param category_name formData string true "分类名称"
// @Success 200 {object} common.CreateProblemResponse "1000 创建成功"
// @Failure 200 {object} common.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} common.CreateProblemResponse "1033 题目分类已存在"
// @Failure 200 {object} common.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/category/update [PUT]
func (a *ApiAdminCategory) UpdateCategory(c *gin.Context) {
	categoryID := c.PostForm("category_id")
	categoryName := c.PostForm("category_name")

	if categoryName == "" || categoryID == "" {
		zap.L().Error("categoryID  or categoryName is null")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	ret := AdminService.UpdateCategory(categoryID, categoryName)
	switch ret.Code {
	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.CategoryIDDoNotExist:
		response.ResponseError(c, response.CodeCategoryIDNotExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}

// DeleteCategory 删除分类
// @Tags Admin API
// @Summary 删除分类
// @Description 删除分类接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param category_id query string true "分类ID"
// @Success 200 {object} common.CreateProblemResponse "1000 删除成功"
// @Failure 200 {object} common.CreateProblemResponse "1001 参数错误"
// @Failure 200 {object} common.CreateProblemResponse "1014 服务器内部错误"
// @Router /admin/category/delete [DELETE]
func (a *ApiAdminCategory) DeleteCategory(c *gin.Context) {
	categoryID := c.Query("category_id")

	if categoryID == "" {
		zap.L().Error("categoryID is null")
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}

	ret := AdminService.DeleteCategory(categoryID)
	switch ret.Code {

	case resp_code.Success:
		response.ResponseSuccess(c, response.CodeSuccess)

	case resp_code.CategoryIsNotEmpty:
		response.ResponseError(c, response.CodeCategoryIsNotEmpty)

	case resp_code.CategoryIDDoNotExist:
		response.ResponseError(c, response.CodeCategoryIDNotExist)

	default:
		response.ResponseError(c, response.CodeInternalServerError)
	}
}
