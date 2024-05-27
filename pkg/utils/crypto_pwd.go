package utils

import (
	"crypto/sha512"
	"encoding/hex"
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

// CheckPwd 验证密码是否正确
func CheckPwd(plainText, cipher string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plainText))
	return err == nil
}

// CryptoSecret 密钥加密
func CryptoSecret(secret string) string {
	data := []byte(secret) // 要加密的数据

	hasher := sha512.New()       // 创建一个 SHA-512 哈希算法实例
	hasher.Write(data)           // 将数据写入哈希算法实例
	encrypted := hasher.Sum(nil) // 计算哈希值

	// 将 SHA-512 哈希值转换为字符串
	encryptedStr := hex.EncodeToString(encrypted)
	return encryptedStr
}
