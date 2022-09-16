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

type RegisterReq struct {
	Phone           string `json:"phone" binding:"required"`
	Code            string `json:"code" binding:"min=4,max=8"`
	Password        string `json:"password" gorm:"-;size:32"`
	ConfirmPassword string `json:"confirm_password" gorm:"-;size:32"`
}

type RegisterResp struct {
	UserId uint64 `json:"user_id"`
	Token  string `json:"token"`
}

type SendPhoneCodeReq struct {
	Phone string `json:"phone" binding:"required"`
	Scene string `json:"scene" binding:"required"`
}
