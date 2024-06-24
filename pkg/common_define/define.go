package common_define

import "fmt"

type CacheKeyMap struct {
	ProblemListPrefix   string
	ProblemDetailPrefix string
	ProblemSearchPrefix string
	HotSearchPrefix     string
	RecentSearchPrefix  string
}

var GlobalCacheKeyMap = CacheKeyMap{
	ProblemListPrefix:   "problem_list",
	ProblemDetailPrefix: "problem_detail",
	ProblemSearchPrefix: "problem_search",
	HotSearchPrefix:     "hot_search",
	RecentSearchPrefix:  "recent_search",
}

var (
	ErrSearchProblem         = fmt.Errorf("problem title not found")
	ErrorBloomFilterNotFound = fmt.Errorf("problem list not found")
)
