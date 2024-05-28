package resp_code

// 服务端请求错误
const (
	Success = 5000 + iota
	SearchDBError
	DBDeleteError
	GenerateNodeError
	EncryptPwdError
	InsertNewUserError
	GenerateTokenError
	SendCodeError
	StoreVerCodeError
	CreateProblemError
	InternalServerError
	GetUserRankError
	Send2MQError
	JSONMarshalError
	RecvFromMQError
	InsertToJudgementError
	RemoveTestFileError
	ReadTestFileError
)
