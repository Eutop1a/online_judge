package verify

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type ApiVerify struct{}

func (v *ApiVerify) InitApiVerify(Router *gin.RouterGroup) {
	verifyApi := v1.ApiGroupApp.ApiVerify

	Router.POST("/check-picture-code", verifyApi.CheckPictureCode)
	Router.POST("/send-email-code", verifyApi.SendEmailCode)
	Router.POST("/send-picture-code", verifyApi.SendPictureCode)
}
