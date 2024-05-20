package my_captcha

import (
	"github.com/mojocn/base64Captcha"
)

func GenerateCaptcha() (id, b64s, ans string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	// 生成验证码图片和相关信息
	id, b64s, ans, err = captcha.Generate()

	if err != nil {
		return "", "", "", err
	}
	return
}
