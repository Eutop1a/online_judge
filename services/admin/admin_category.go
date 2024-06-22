package admin

import (
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/common/response"
)

type AdminCategoryService struct{}

// AddCategory 添加分类
func (a *AdminCategoryService) AddCategory(categoryName string) (resp response.Response) {
	// 添加分类，categoryName 冲突的情况 mysql 会报错 ErrCategoryAlreadyExist
	err := mysql.CreateNewCategory(categoryName)
	if err != nil {
		if err == mysql.ErrCategoryAlreadyExist {
			resp.Code = resp_code.CategoryTypeAlreadyExist
			zap.L().Error("service-AddCategory-CreateNewCategory ", zap.Error(err))
			return
		}
		resp.Code = resp_code.InternalServerError
		zap.L().Error("service-AddCategory-CreateNewCategory ", zap.Error(err))
		return
	}
	resp.Code = resp_code.Success
	return
}

// UpdateCategory 更新分类
func (a *AdminCategoryService) UpdateCategory(categoryID, categoryName string) (resp response.Response) {
	// 添加分类，categoryName 冲突的情况 mysql 会报错 ErrCategoryAlreadyExist
	err := mysql.UpdateCategoryDetail(categoryID, categoryName)
	if err != nil {
		if err == mysql.ErrCategoryNotFound {
			resp.Code = resp_code.CategoryIDDoNotExist
			zap.L().Error("service-UpdateCategory-UpdateCategory ",
				zap.String("categoryID", categoryID),
				zap.Error(err))
			return
		}
		zap.L().Error("service-UpdateCategory-UpdateCategory ",
			zap.String("category_id", categoryID),
			zap.Error(err))
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
