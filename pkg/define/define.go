package define

import "fmt"

type CacheKeyMap struct {
	ProblemListPrefix   string
	ProblemDetailPrefix string
	ProblemSearchPrefix string
	HotSearchPrefix     string
	RecentSearchPrefix  string
	ProblemIDListPrefix string
}

var GlobalCacheKeyMap = CacheKeyMap{
	ProblemListPrefix:   "problem_list",
	ProblemDetailPrefix: "problem_detail",
	ProblemSearchPrefix: "problem_search",
	HotSearchPrefix:     "hot_search",
	RecentSearchPrefix:  "recent_search",
	ProblemIDListPrefix: "problem_id_list",
}

var (
	ErrSearchProblem       = fmt.Errorf("problem title not found")
	ErrBloomFilterNotFound = fmt.Errorf("problem list not found")
	ErrProblemIDNotFound   = fmt.Errorf("problem id not found")
	ErrNoProblemIDFound    = fmt.Errorf("no problem ID found in cache")
)
