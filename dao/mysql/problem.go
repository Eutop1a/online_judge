package mysql

import "gorm.io/gorm"

// GetProblemList 获取题目列表
func GetProblemList() (*[]Problems, error) {
	var problemList []Problems

	// 执行查询并只选取题号、题目和难度字段
	err := DB.Model(&Problems{}).Select("problem_id, title, difficulty").
		Find(&problemList).Error
	if err != nil {
		// 处理错误
		return nil, err
	}
	return &problemList, nil
}

// GetProblemDetail 获取单个题目详细信息
func GetProblemDetail(pid string) (problem *Problems, err error) {
	//err = DB.Preload("TestCases").Where("problem_id = ?", pid).First(&problem).Error
	err = DB.Where("problem_id = ?", pid).Preload("TestCases", func(db *gorm.DB) *gorm.DB {
		return db.Limit(2) // 在这里使用 Limit 方法限制 TestCases 的数量
	}).First(&problem).Error

	if err != nil {
		return nil, err
	}
	return
}

// GetEntireProblem 获取题目的全部信息
func GetEntireProblem(pid string) (problem *Problems, err error) {
	err = DB.Where("problem_id = ?", pid).Preload("TestCases").First(&problem).Error
	if err != nil {
		return nil, err
	}
	return
}

// CreateProblem 创建题目
func CreateProblem(problem *Problems) error {
	return DB.Create(problem).Error
}

// UpdateProblem 更新题目
func UpdateProblem(problem *Problems) error {
	// TODO: 开启事务
	// TODO: 更新问题基础信息
	// TODO: 更新关联的问题测试样例
	ts := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			ts.Rollback()
		}
	}()

	err := ts.Where("problem_id = ?", problem.ProblemID).Updates(&problem).Error
	if err != nil {
		return err
	}
	if len(problem.TestCases) != 0 {
		ts.Where("pid = ?", problem.ProblemID).Delete(new(TestCase))
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
	// TODO: 开启事务，删除题目基本信息，删除题目测试样例
	ts := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			ts.Rollback()
		}
	}()

	err := ts.Where("problem_id = ?", pid).Delete(&Problems{}).Error
	if err != nil {
		return err
	}

	err = ts.Where("pid = ?", pid).Delete(&TestCase{}).Error
	if err != nil {
		return err
	}
	ts.Commit()
	return nil
}

// CheckProblemTitle 检查题目标题是否已经存在
func CheckProblemTitle(title string, num *int64) error {
	return DB.Model(&Problems{}).Where("title = ?", title).Count(num).Error
}

// CheckProblemID 检查题目id是否已经存在
func CheckProblemID(id string, num *int64) error {
	return DB.Model(&Problems{}).Where("problem_id = ?", id).Count(num).Error
}

func GetProblemID(title string) (problemID string, err error) {
	var problem Problems
	err = DB.Model(&Problems{}).Where("title = ?", title).First(&problem).Error
	if err != nil {
		return "", err
	}
	return problem.ProblemID, nil
}
