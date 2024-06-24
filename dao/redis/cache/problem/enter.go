package problem

import "online_judge/services"

type CacheGroup struct {
	CacheProblem
}

var (
	ProblemService = services.ServiceGroupApp.ProblemService
)
