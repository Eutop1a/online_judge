package controller

type _ResponseRegister struct {
	Token string `json:"token"`
	Error error  `json:"error"`
	Msg   error  `json:"msg"`
}

type _ResponseSendCode struct {
	Msg string `json:"msg"`
}
