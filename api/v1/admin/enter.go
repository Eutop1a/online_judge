package admin

import "online_judge/services"

type ApiGroup struct {
	ApiAdminUser
	ApiAdminProblem
	ApiAdminCategory
}

var (
	AdminService = services.ServiceGroupApp.AdminService
)
