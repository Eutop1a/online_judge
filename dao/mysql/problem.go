package mysql

func GetProblemList() ([]Problems, error) {
	var problemList []Problems

	// 执行查询并只选取题号、题目和难度字段
	err := DB.Select("problemID, title, difficulty").
		Find(&problemList).Error

	if err != nil {
		// 处理错误
		return nil, err
	}

	return problemList, nil
}

func CreateProblem(pid, title, content, difficulty string, runtime, memory int, testCase []*TestCase) error {
	var problemDetail Problems

	problemDetail.ProblemID = pid
	problemDetail.Title = title
	problemDetail.Content = content
	problemDetail.Difficulty = difficulty
	problemDetail.MaxRuntime = runtime
	problemDetail.MaxMemory = memory
	problemDetail.TestCases = testCase

	return DB.Create(&problemDetail).Error
}

func CheckProblemTitle(title string, num *int64) error {
	return DB.Model(&Problems{}).Where("title = ?", title).Count(num).Error
}
