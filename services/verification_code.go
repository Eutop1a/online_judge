package services

import (
	"fmt"
	"go.uber.org/zap"
	"online-judge/dao/redis"
	mycaptcha "online-judge/pkg/my_captcha"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	"time"
)

// SendEmailCode 发送验证码接口
func SendEmailCode(userEmail string) (resCode int) {
	// 创建验证码
	code, ts := utils.CreateVerificationCode()
	// 发送验证码
	err := utils.SendCode(userEmail, code)
	if err != nil {
		resCode = resp.SendCodeError
		zap.L().Error("services-SendEmailCode-SendCode ", zap.Error(err))
		return
	}
	// redis保存email和验证码的键值对
	err = redis.StoreVerificationCode(userEmail, code, ts)
	if err != nil {
		resCode = resp.StoreVerCodeError
		zap.L().Error("services-SendEmailCode-StoreVerificationCode ", zap.Error(err))
		return
	}

	resCode = resp.Success
	return
}

// SendPictureCode 发送图片验证码
func SendPictureCode(username string) (pic string, err error) {
	// 单例模式的验证码实例
	_, b64s, ans, err := mycaptcha.GenerateCaptcha()

	if err != nil {
		zap.L().Error("services-SendCode-GenerateCaptcha ", zap.Error(err))
		return "", err
	}
	// 获取当前时间
	ts := time.Now().Unix()
	err = redis.StorePictureCode(username, ans, ts)
	if err != nil {
		zap.L().Error("services-SendCode-StorePictureCode ", zap.Error(err))
		return "", err
	}
	return b64s, nil
}

// CheckCode 检查图片验证码
func CheckCode(username, code string) (bool, error) {
	ans, err := redis.GetPictureCode(username)
	if err != nil {
		zap.L().Error("services-CheckCode-GetPictureCode "+
			fmt.Sprintf("do not have this username %s ", username), zap.Error(err))
		return true, err
	}
	if ans != code {
		zap.L().Error("services-CheckCode-GetPictureCode wrong picture code from: " + username)
		return false, nil
	}

	return true, nil
}
