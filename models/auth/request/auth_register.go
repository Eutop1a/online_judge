package request

type AuthRegisterReq struct {
	UserID   int64  `form:"user_id" json:"user_id"`
	UserName string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Email    string `form:"email" json:"email" validate:"required"`
	Code     string `form:"code" json:"code" validate:"required"`
	//Role             bool      `form:"role" json:"role"`
	// true is Admin, false is user
}
