package request

type AdminDeleteUserReq struct {
	UserID int64 `json:"user_id" form:"user_id"`
}

type AdminAddAdminReq struct {
	Username string `json:"username" form:"username"`
}

type AdminAddSuperAdminReq struct {
	Username string `json:"username" form:"username"`
	Secret   string `json:"secret" form:"secret"`
}
