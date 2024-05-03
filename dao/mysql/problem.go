package mysql

import "gorm.io/gorm"

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

func CreateProblem(problemDetail *Problems) error {
	return DB.Create(problemDetail).Error
}

func CheckProblemTitle(title string, num *int64) error {
	return DB.Model(&Problems{}).Where("title = ?", title).Count(num).Error
}
