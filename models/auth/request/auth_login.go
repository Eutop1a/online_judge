package request

type AuthLoginReq struct {
	UserName string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}
