package verify

import (
	"online_judge/services"
)

type ApiGroup struct {
	ApiVerify
}

var (
	VerifyService = services.ServiceGroupApp.VerifyService
)
