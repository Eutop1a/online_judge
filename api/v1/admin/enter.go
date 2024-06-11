package admin

import "online_judge/services"

type ApiGroup struct {
	ApiAdmin
}

var (
	AdminService = services.ServiceGroupApp.AdminService
)
