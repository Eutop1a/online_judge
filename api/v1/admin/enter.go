package admin

import (
	"online_judge/dao/redis/cache"
	"online_judge/services"
)

type ApiGroup struct {
	ApiAdminUser
	ApiAdminProblem
	ApiAdminCategory
}

var (
	AdminService = services.ServiceGroupApp.AdminService
	CacheService = cache.CacheGroupApp.CacheAdmin
)
