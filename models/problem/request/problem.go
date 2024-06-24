package request

type GetProblemListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetProblemDetailReq struct {
	ProblemID string `json:"problem_id" form:"problem_id"`
}

type GetProblemIDReq struct {
	Title string `json:"title" form:"title"`
}

type GetProblemRandomReq struct {
}

type SearchProblemReq struct {
	Title string `json:"msg" form:"msg"`
}
