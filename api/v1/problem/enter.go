package problem

import "online_judge/services"

type ApiGroup struct {
	ApiProblem
}

var (
	ProblemService = services.ServiceGroupApp.ProblemService
)
