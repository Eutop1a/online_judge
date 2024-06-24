package response

type SearchProblemResp struct {
	Data []*ProblemResponse `json:"data"`
}

type ProblemResponse struct {
	ID         int                `json:"id"`
	ProblemID  string             `json:"problem_id"`
	Title      string             `json:"title"`
	Difficulty string             `json:"difficulty"`
	Categories []CategoryResponse `json:"categories"`
	TestCases  []TestCaseResponse `json:"test_cases"`
}

type ProblemDetailResponse struct {
	ID         int                `json:"id"`
	ProblemID  string             `json:"problem_id"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Difficulty string             `json:"difficulty"`
	MaxRuntime int                `json:"max_runtime"` // 时间限制
	MaxMemory  int                `json:"max_memory"`  // 内存限制
	Categories []CategoryResponse `json:"categories"`
	TestCases  []TestCaseResponse `json:"test_cases"`
}

type CategoryResponse struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	ParentName string `json:"parent_name"`
}

type TestCaseResponse struct {
	TID      string `json:"tid"`
	PID      string `json:"pid"`
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type GetProblemListResp struct {
	Data  []*ProblemResponse `json:"data"`
	Count int64              `json:"count"`
	Size  int                `json:"size"`
	Page  int                `json:"page"`
}
