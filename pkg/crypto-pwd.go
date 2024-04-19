package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

// CryptoPwd 密码加密
func CryptoPwd(pwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// DecryptPwd 验证密码是否正确
func DecryptPwd(cipher, plainText string) bool {
	//fmt.Println(cipher)
	//fmt.Println(plainText)
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plainText))
	if err != nil {
		return false
	}
	return true
}
