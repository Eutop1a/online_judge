package verify

import (
	"fmt"
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/redis/cache/verify"
	"online_judge/models/verify/request"
	mycaptcha "online_judge/pkg/my_captcha"
	"online_judge/pkg/utils"
	"time"
)

type VerifyService struct{}

// SendEmailCode 发送验证码接口
func (v *VerifyService) SendEmailCode(request request.SendEmailCodeReq) (resCode int) {
	// 创建验证码
	code, ts := utils.CreateVerificationCode()
	// 发送验证码
	err := utils.SendCode(request.UserEmail, code)
	if err != nil {
		resCode = resp_code.SendCodeError
		zap.L().Error("services-SendEmailCode-SendCode ", zap.Error(err))
		return
	}
	// redis保存email和验证码的键值对
	err = verify.StoreVerifyCode(request.UserEmail, code, ts)
	if err != nil {
		resCode = resp_code.StoreVerCodeError
		zap.L().Error("services-SendEmailCode-StoreVerifyCode ", zap.Error(err))
		return
	}

	resCode = resp_code.Success
	return
}

// SendPictureCode 发送图片验证码
func (v *VerifyService) SendPictureCode(request request.SendPictureCodeReq) (pic string, err error) {
	// 单例模式的验证码实例
	_, b64s, ans, err := mycaptcha.GenerateCaptcha()

	if err != nil {
		zap.L().Error("services-SendCode-GenerateCaptcha ", zap.Error(err))
		return "", err
	}
	// 获取当前时间
	ts := time.Now().Unix()
	err = verify.StorePictureCode(request.Username, ans, ts)
	if err != nil {
		zap.L().Error("services-SendCode-StorePictureCode ", zap.Error(err))
		return "", err
	}
	return b64s, nil
}

// CheckCode 检查图片验证码
func (v *VerifyService) CheckCode(request request.CheckCodeReq) (bool, error) {
	ans, err := verify.GetPictureCode(request.Username)
	if err != nil {
		zap.L().Error("services-CheckCode-GetPictureCode "+
			fmt.Sprintf("do not have this username %s ", request.Username), zap.Error(err))
		return true, err
	}
	if ans != request.Code {
		zap.L().Error("services-CheckCode-GetPictureCode wrong picture code from: " + request.Username)
		return false, nil
	}

	return true, nil
}
