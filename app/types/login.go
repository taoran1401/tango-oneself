package types

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" gorm:"-;size:32"`
	Code     string `json:"code" binding:"min=4,max=8"`
}

type LoginResp struct {
	UserId uint64 `json:"user_id"`
	Token  string `json:"token"`
}

type SendPhoneCodeReq struct {
	Phone string `json:"phone" binding:"required"`
}
