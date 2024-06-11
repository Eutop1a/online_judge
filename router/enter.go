package router

import (
	"online_judge/router/admin"
	"online_judge/router/auth"
	"online_judge/router/evaluation"
	"online_judge/router/leaderboard"
	"online_judge/router/problem"
	"online_judge/router/submission"
	"online_judge/router/user"
	"online_judge/router/verify"
)

type RouterGroup struct {
	Admin       admin.RouterGroup
	User        user.RouterGroup
	Auth        auth.RouterGroup
	Verify      verify.RouterGroup
	Problem     problem.RouterGroup
	Submission  submission.RouterGroup
	Leaderboard leaderboard.RouterGroup
	Evaluation  evaluation.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
