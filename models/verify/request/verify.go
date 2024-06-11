package request

type SendEmailCodeReq struct {
	UserEmail string `json:"user_email" form:"user_email"`
}

type SendPictureCodeReq struct {
	Username string `json:"username" form:"username"`
}

type CheckCodeReq struct {
	Username string `json:"username" form:"username"`
	Code     string `json:"code" form:"code"`
}
