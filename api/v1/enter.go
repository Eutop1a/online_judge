package v1

import (
	"online_judge/api/v1/admin"
	"online_judge/api/v1/auth"
	"online_judge/api/v1/evaluation"
	"online_judge/api/v1/leaderboard"
	"online_judge/api/v1/problem"
	"online_judge/api/v1/submission"
	"online_judge/api/v1/user"
	"online_judge/api/v1/verify"
)

type ApiGroup struct {
	ApiAdmin       admin.ApiGroup
	ApiUser        user.ApiGroup
	ApiAuth        auth.ApiGroup
	ApiVerify      verify.ApiGroup
	ApiProblem     problem.ApiGroup
	ApiSubmission  submission.ApiGroup
	ApiLeaderboard leaderboard.ApiGroup
	ApiEvaluation  evaluation.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
