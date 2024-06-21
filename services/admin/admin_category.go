package admin

import (
	"fmt"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/common/response"
)

type AdminCategoryService struct{}

// AddCategory 添加分类
func (a *AdminCategoryService) AddCategory(categoryName string) (resp response.Response) {
	// 检查是否已经存在该类型
	ok, err := mysql.CheckCategoryByName(categoryName)
	if err != nil {
		resp.Code = resp_code.InternalServerError
		zap.L().Error("check category name error ", zap.Error(err))
		return
	}
	if ok {
		// 该类型已经存在
		resp.Code = resp_code.CategoryTypeAlreadyExist
		zap.L().Error(fmt.Sprintf("category %s already exist ", categoryName), zap.Error(err))
		return
	}
	// 如果不存在则添加
	err = mysql.CreateNewCategory(categoryName)
	if err != nil {
		resp.Code = resp_code.InternalServerError
		zap.L().Error("add new category error ", zap.Error(err))
		return
	}
	resp.Code = resp_code.Success
	return
}

// UpdateCategory 更新分类
func (a *AdminCategoryService) UpdateCategory(categoryID, categoryName string) (resp response.Response) {
	// 检查要修改的categoryID是否存在
	ok, err := mysql.CheckCategoryById(categoryID)
	if err != nil {
		resp.Code = resp_code.InternalServerError
		zap.L().Error("check category name error ", zap.Error(err))
		return
	}
	if !ok {
		// 该类型已经存在
		resp.Code = resp_code.CategoryIDDoNotExist
		zap.L().Error(fmt.Sprintf("category %s already exist ", categoryID), zap.Error(err))
		return
	}

	err = mysql.UpdateCategoryDetail(categoryID, categoryName)
	if err != nil {
		zap.L().Error("update category error ", zap.String("category_id", categoryID), zap.Error(err))
		resp.Code = resp_code.InternalServerError
		return
	}

	resp.Code = resp_code.Success
	return
}

// DeleteCategory 删除分类
func (a *AdminCategoryService) DeleteCategory(categoryID string) (resp response.Response) {
	ok, err := mysql.DeleteCategoryById(categoryID)
	if err != nil {
		resp.Code = resp_code.InternalServerError
		zap.L().Error("delete category error ", zap.Error(err))
		return
	}
	if !ok {
		resp.Code = resp_code.CategoryIsNotEmpty
		zap.L().Error("there are still problems under the category ", zap.Error(err))
		return
	}
	resp.Code = resp_code.Success
	return
}
