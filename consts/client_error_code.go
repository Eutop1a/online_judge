package consts

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// 客户端请求错误
const (
	EmailAlreadyExist = 4000 + iota
	UsernameAlreadyExist
	InvalidateEmailFormat
	ErrorVerCode
	ExpiredVerCode
	NotExistUsername
	NotExistUserID
	ErrorPwd
	ProblemAlreadyExist
	ProblemNotExist
	UnsupportedLanguage
	SecretError
	UserIDAlreadyExist
	//ErrInvalidCredentials
)
