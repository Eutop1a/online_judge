package utils

import (
	"regexp"
)

var emailRegex *regexp.Regexp

func init() {
	// 预编译正则表达式
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
}

func ValidateEmail(email string) bool {
	// 使用预编译的正则表达式
	return emailRegex.MatchString(email)
}
