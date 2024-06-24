package v1

import (
	"online_judge/api/v1/admin"
	"online_judge/api/v1/auth"
	"online_judge/api/v1/category"
	"online_judge/api/v1/evaluation"
	"online_judge/api/v1/leaderboard"
	"online_judge/api/v1/problem"
	"online_judge/api/v1/submission"
	"online_judge/api/v1/user"
	"online_judge/api/v1/verify"
)

type ApiGroupInterface interface {
	GetAdminApiGroup() admin.ApiGroup
	GetUserApiGroup() user.ApiGroup
	GetAuthApiGroup() auth.ApiGroup
	GetVerifyApiGroup() verify.ApiGroup
	GetProblemApiGroup() problem.ApiGroup
	GetSubmissionApiGroup() submission.ApiGroup
	GetLeaderboardApiGroup() leaderboard.ApiGroup
	GetEvaluationApiGroup() evaluation.ApiGroup
	GetCategoryApiGroup() category.ApiGroup
}

type ApiGroup struct {
	ApiAdmin       admin.ApiGroup
	ApiUser        user.ApiGroup
	ApiAuth        auth.ApiGroup
	ApiVerify      verify.ApiGroup
	ApiProblem     problem.ApiGroup
	ApiSubmission  submission.ApiGroup
	ApiLeaderboard leaderboard.ApiGroup
	ApiEvaluation  evaluation.ApiGroup
	ApiCategory    category.ApiGroup
}

//func (a *ApiGroup) GetAdminApiGroup() admin.ApiGroup {
//	return a.ApiAdmin
//}
//
//func (a *ApiGroup) GetUserApiGroup() user.ApiGroup {
//	return a.ApiUser
//}
//
//func (a *ApiGroup) GetAuthApiGroup() auth.ApiGroup {
//	return a.ApiAuth
//}
//
//func (a *ApiGroup) GetVerifyApiGroup() verify.ApiGroup {
//	return a.ApiVerify
//}
//
//func (a *ApiGroup) GetProblemApiGroup() problem.ApiGroup {
//	return a.ApiProblem
//}
//
//func (a *ApiGroup) GetSubmissionApiGroup() submission.ApiGroup {
//	return a.ApiSubmission
//}
//
//func (a *ApiGroup) GetLeaderboardApiGroup() leaderboard.ApiGroup {
//	return a.ApiLeaderboard
//}
//
//func (a *ApiGroup) GetEvaluationApiGroup() evaluation.ApiGroup {
//	return a.ApiEvaluation
//}
//
//func (a *ApiGroup) GetCategoryApiGroup() category.ApiGroup {
//	return a.ApiCategory
//}
//
//var ApiGroupApp ApiGroupInterface = &ApiGroup{}

var ApiGroupApp = new(ApiGroup)
