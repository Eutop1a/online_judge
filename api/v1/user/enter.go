package user

import "online_judge/services"

type ApiGroup struct {
	ApiUser
}

var (
	UserService = services.ServiceGroupApp.UserService
)
