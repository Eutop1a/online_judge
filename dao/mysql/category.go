package mysql

import (
	"github.com/go-sql-driver/mysql"
	"online_judge/pkg/utils"
)

// CheckCategoryByName 通过 name 检查是否存在这个 category
func CheckCategoryByName(categoryName string) (bool, error) {
	var count int64
	err := DB.Model(&Category{}).Where("name = ?", categoryName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckCategoryById 通过 id 检查是否存在这个 category
func CheckCategoryById(categoryID string) (bool, error) {
	var count int64
	err := DB.Model(&Category{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateNewCategory 创建新的category
func CreateNewCategory(categoryName string) error {
	err := DB.Create(&Category{
		Name:       categoryName,
		CategoryID: utils.GetUUID(),
	}).Error

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return ErrCategoryAlreadyExist
			}
		}
		return err
	}
	return nil
}

// UpdateCategoryDetail 更新 category name
func UpdateCategoryDetail(categoryID, categoryName string) error {
	result := DB.Model(&Category{}).Where("category_id = ?", categoryID).
		Updates(map[string]interface{}{"name": categoryName})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// DeleteCategoryById 删除 category
func DeleteCategoryById(categoryID string) (bool, error) {
	var count int64
	// 搜索分类下是否有题目，使用join连接两张表
	err := DB.Model(&ProblemCategory{}).
		Joins("JOIN category ON category.id = problem_category.CategoryId AND category.identity = ?",
			categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	// 分类下面还有题目，禁止删除
	if count > 0 {
		return false, nil
	}
	err = DB.Model(&Category{}).Where("category_id = ?", categoryID).Delete(&Category{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetCategoryID 获取分类ID
func GetCategoryID(categoryName string) (string, error) {
	var tmp Category
	err := DB.Model(&Category{}).Where("name = ?", categoryName).First(&tmp).Error
	if err != nil {
		return "", err
	}
	return tmp.CategoryID, nil
}

func CreateProblemCategory(problemID, categoryID string) error {
	var tmp ProblemCategory
	tmp.ProblemIdentity = problemID
	tmp.CategoryIdentity = categoryID
	return DB.Create(&tmp).Error
}

func DeleteProblemCategoryByProblemID(problemID string) error {
	var tmp ProblemCategory
	tmp.ProblemIdentity = problemID
	return DB.Delete(&tmp).Error
}

func UpdateProblemCategoryByProblemID(problemID, categoryID string) error {
	var tmp ProblemCategory
	tmp.ProblemIdentity = problemID
	tmp.CategoryIdentity = categoryID

	return DB.Save(&tmp).Error
}
