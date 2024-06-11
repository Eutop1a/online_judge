package response

import "online_judge/dao/mysql"

type GetProblemListResp struct {
	Data  []*mysql.Problems
	Count int64
	Size  int
	Page  int
}
