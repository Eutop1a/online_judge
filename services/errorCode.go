package services

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
	UpdateLoginDataError
)

// 服务端请求错误
const (
	Success = 5000 + iota
	SearchDBError
	DBSaveError
	DBDeleteError
	GenerateNodeError
	EncryptPwdError
	InsertNewUserError
	GenerateTokenError
	SendCodeError
	StoreVerCodeError
)
