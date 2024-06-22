package mysql

import (
	"gorm.io/gorm"
	"math/rand"
)

// GetProblemList 获取题目列表
func GetProblemList(page, size int, count *int64) ([]Problems, error) {
	offset := (page - 1) * size // 从哪里开始查询，例如page = 1，应该从数据库的第0条记录开始查询

	var problemList []Problems

	err := DB.Model(&Problems{}).Count(count).
		Preload("ProblemCategories.Category").
		Select("id, problem_id, title, difficulty").
		Offset(offset).Limit(size).Find(&problemList).Error
	if err != nil {
		// 处理错误
		return nil, err
	}

	return problemList, nil
}

// GetProblemDetail 获取单个题目详细信息
func GetProblemDetail(pid string) (problem *Problems, err error) {
	err = DB.Where("problem_id = ?", pid).
		Preload("TestCases", func(db *gorm.DB) *gorm.DB {
			return db.Limit(2) // 在这里使用 Limit 方法限制 TestCases 的数量
		}).
		Preload("ProblemCategories.Category").
		First(&problem).Error

	if err != nil {
		return nil, err
	}
	return
}

// GetProblemRandom 随机获取一个题目
func GetProblemRandom() (problem *Problems, err error) {
	var problemsList []Problems
	err = DB.Model(&Problems{}).Find(&problemsList).Error
	if err != nil {
		return nil, err
	}
	//fmt.Println("len: ", len(problemsList))

	randomIdx := rand.Intn(len(problemsList))
	problemIdx := &problemsList[randomIdx].ID

	err = DB.Where("id = ?", problemIdx).Preload("TestCases", func(db *gorm.DB) *gorm.DB {
		return db.Limit(2)
	}).First(&problem).Error

	if err != nil {
		return nil, err
	}
	return problem, nil
}

// GetEntireProblem 获取题目的全部信息
func GetEntireProblem(pid string) (problem *Problems, err error) {
	result := DB.Where("problem_id = ?", pid).Preload("TestCases").First(&problem)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return problem, nil
}

// CreateProblem 创建题目
func CreateProblem(problem *Problems) error {
	return DB.Create(problem).Error
}

// UpdateProblem 更新题目
func UpdateProblem(problem *Problems, oldProblemID string, category []string) error {
	// TODO:更新问题基础信息
	// TODO:更新关联的问题分类
	// TODO:更新关联的问题测试样例
	// 开启事务
	ts := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			ts.Rollback()
		}
	}()

	// TODO:更新问题基础信息
	err := ts.Where("problem_id = ?", problem.ProblemID).Updates(&problem).Error
	if err != nil {
		return err
	}

	//TODO:更新关联的问题分类
	//浅复制即可
	/*先删除原来的再添加,若categoryIds为空表示无需修改分类*/
	if len(category) != 0 {
		err = ts.Model(&ProblemCategory{}).Where("problem_id = ?", oldProblemID).Delete(&ProblemCategory{}).Error
		if err != nil {
			return err
		}
		pcSlice := make([]*ProblemCategory, 0)
		for _, v := range category {
			pcSlice = append(pcSlice, &ProblemCategory{
				ProblemIdentity:  oldProblemID,
				CategoryIdentity: v,
			})
		}
		err = ts.Create(&pcSlice).Error
		if err != nil {
			return err
		}
	}

	//TODO:更新关联的问题测试样例
	/*先删除原来的再添加,若testCases为空表示无需修改测试用例*/
	if len(problem.TestCases) != 0 {
		ts.Where("pid = ?", problem.ProblemID).Delete(&TestCase{})
		err = ts.Create(&problem.TestCases).Error
		if err != nil {
			return err
		}
	}
	ts.Commit()
	return nil
}

// DeleteProblem 删除题目
func DeleteProblem(pid string) error {
	//TODO: 开启事务处理
	ts := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			ts.Rollback()
		}
	}()

	//TODO:删除问题本体
	result := ts.Where("problem_id = ?", pid).Delete(&Problems{})

	if result.Error != nil {
		return result.Error
	}
	// 没有收到影响的行，说明pid不存在，返回错误
	if result.RowsAffected == 0 {
		return ErrProblemIDNotExist
	}

	//TODO:删除问题分类关联表的对应内容
	err := ts.Where("problem_id = ?", pid).Delete(&ProblemCategory{}).Error
	if err != nil {
		return err
	}

	//TODO:删除测试用例
	err = ts.Where("pid = ?", pid).Delete(&TestCase{}).Error
	if err != nil {
		return err
	}
	ts.Commit()

	return nil
}

// CheckProblemTitleExists 检查题目标题是否已经存在
func CheckProblemTitleExists(title string) (bool, error) {
	var count int64
	err := DB.Model(&Problems{}).Where("title = ?", title).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckProblemIDExists 检查题目id是否已经存在
func CheckProblemIDExists(id string) (bool, error) {
	var count int64
	err := DB.Model(&Problems{}).Where("problem_id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetProblemID 获取题目ID
func GetProblemID(title string) (problemID string, err error) {
	var problem Problems
	err = DB.Model(&Problems{}).Where("title = ?", title).First(&problem).Error
	if err != nil {
		return "", err
	}
	return problem.ProblemID, nil
}

// CreateProblemWithFile 创建以文件为输入输出的题目
func CreateProblemWithFile(problem *ProblemWithFile) error {
	return DB.Model(&ProblemWithFile{}).Create(problem).Error
}

// DeleteProblemWithFile 删除题目
func DeleteProblemWithFile(pid string) (string, error) {
	var tmp ProblemWithFile
	if err := DB.Model(&ProblemWithFile{}).Where("problem_id = ?", pid).Find(&tmp).Error; err != nil {
		return "", err
	}
	err := DB.Model(&ProblemWithFile{}).Where("problem_id = ?", pid).Delete(&ProblemWithFile{}).Error
	if err != nil {
		return "", err
	}

	return tmp.ProblemID, nil
}

// CheckProblemIDWithFile 检查题目id是否已经存在
func CheckProblemIDWithFile(id string, num *int64) error {
	return DB.Model(&ProblemWithFile{}).Where("problem_id = ?", id).Count(num).Error
}

// UpdateProblemWithFile 更新题目
func UpdateProblemWithFile(problem *ProblemWithFile) error {
	return DB.Model(&ProblemWithFile{}).Where("problem_id = ?", problem.ProblemID).Updates(&problem).Error
}

// CheckProblemTitleWithFile 检查题目标题是否已经存在
func CheckProblemTitleWithFile(title string, num *int64) error {
	return DB.Model(&ProblemWithFile{}).Where("title = ?", title).Count(num).Error
}

// GetEntireProblemWithFile 获取题目的全部信息
func GetEntireProblemWithFile(pid string) (problem *ProblemWithFile, err error) {
	err = DB.Model(&ProblemWithFile{}).Where("problem_id = ?", pid).First(&problem).Error
	if err != nil {
		return nil, err
	}
	return
}

// SearchProblemByMsg 模糊搜索题目获取题目列表
func SearchProblemByMsg(problems *[]Problems, searchQuery string) error {
	return DB.Model(&Problems{}).Where("title LIKE ? OR content LIKE ?", searchQuery, searchQuery).
		Preload("ProblemCategories.Category").
		Find(&problems).Error
	//return DB.Model(&Problems{}).Select("id, problem_id, title, difficulty").
	//	Where("title LIKE ?", searchQuery).Order("id asc").Find(&problems).Error
}

// GetProblemListByCategory 根据分类ID获取题目列表
func GetProblemListByCategory(categoryID string) ([]Problems, error) {
	var problems []Problems

	err := DB.Joins("JOIN problem_category ON problems.problem_id = problem_category.problem_id").
		Where("problem_category.category_id = ?", categoryID).
		Preload("ProblemCategories.Category").
		Find(&problems).Error

	if err != nil {
		return nil, err
	}

	return problems, nil
}

// GetCategoryList 获取分类列表
func GetCategoryList() ([]*Category, error) {
	var tmp []*Category
	err := DB.Model(&Category{}).Find(&tmp).Error
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

// GetProblemsByCategoryName 根据分类名称获取题目列表
func GetProblemsByCategoryName(categoryName string) ([]Problems, error) {
	var problems []Problems

	err := DB.Joins("JOIN problem_category ON problems.problem_id = problem_category.problem_id").
		Joins("JOIN category ON category.category_id = problem_category.category_id").
		Where("category.name = ?", categoryName).
		Preload("ProblemCategories.Category").
		Find(&problems).Error

	if err != nil {
		return nil, err
	}

	return problems, nil
}
