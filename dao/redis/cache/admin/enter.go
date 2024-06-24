package admin

import "online_judge/services"

type CacheGroup struct {
}

var (
	AdminService = services.ServiceGroupApp.AdminService
)
