package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = 12
)

// CryptoPwd 密码加密
func CryptoPwd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DecryptPwd 验证密码是否正确
func DecryptPwd(cipher, plainText string) error {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plainText))
	return err
}
