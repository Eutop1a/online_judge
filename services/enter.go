package services

import (
	"online_judge/services/admin"
	"online_judge/services/auth"
	"online_judge/services/category"
	"online_judge/services/evaluation"
	"online_judge/services/leaderboard"
	"online_judge/services/problem"
	"online_judge/services/submission"
	"online_judge/services/user"
	"online_judge/services/verify"
)

type ServiceGroup struct {
	AdminService       admin.ServiceGroup
	UserService        user.ServiceGroup
	AuthService        auth.ServiceGroup
	VerifyService      verify.ServiceGroup
	ProblemService     problem.ServiceGroup
	SubmissionService  submission.ServiceGroup
	LeaderboardService leaderboard.ServiceGroup
	EvaluationService  evaluation.ServiceGroup
	CategoryService    category.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
