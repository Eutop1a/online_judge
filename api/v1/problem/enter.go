package problem

import (
	"online_judge/dao/redis/cache"
	"online_judge/services"
)

type ApiGroup struct {
	ApiProblem
}

var (
	ProblemService = services.ServiceGroupApp.ProblemService
	ProblemCache   = cache.CacheGroupApp.CacheProblem
)
