package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
	"math/rand"
	"net/smtp"
	"online_judge/dao/redis/cache/verify"
	"strconv"
	"time"
)

func SendCode(toUser, code string) error {
	e := email.NewEmail()
	e.Subject = "Verification code"
	e.HTML = []byte("your verify code is: <b>" + code + "</b>" + "   expired time: " +
		strconv.FormatInt(verify.Expired/60, 10) + " minutes")

	e.From = "Eutop1a <w1905700640@163.com>"
	e.To = []string{toUser}

	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "w1905700640@163.com", "TCJZVOJENDBLMUWC", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: false, ServerName: "smtp.163.com"})
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error sending email to %s:%v", toUser, err))
	}
	return err
}

func CreateVerificationCode() (string, int64) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	t1 := rand.Int() % 1000000
	ret := fmt.Sprintf("%06d", t1)
	timestamp := time.Now().Unix()
	return ret, timestamp

}
