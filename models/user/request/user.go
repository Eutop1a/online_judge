package request

type GetUserDetailReq struct {
	UserID int64 `form:"user_id" json:"user_id"`
}

type UpdateUserDetailReq struct {
	UserID   int64  `form:"user_id" json:"user_id"`
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Email    string `form:"email" json:"email" validate:"required"`
	Code     string `form:"code" json:"code" validate:"required"`
}

type GetUserIDReq struct {
	Username string `form:"username" json:"username" validate:"required"`
}
