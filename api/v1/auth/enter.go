package auth

import "online_judge/services"

type ApiGroup struct {
	ApiAuth
}

var (
	AuthService = services.ServiceGroupApp.AuthService
)
