package define

type CacheKeyMap struct {
	ProblemListPrefix   string
	ProblemDetailPrefix string
}

var GlobalCacheKeyMap = CacheKeyMap{
	ProblemListPrefix:   "problemlist",
	ProblemDetailPrefix: "problem",
}
