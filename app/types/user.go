package types

//用户信息
type UserBase struct {
	Id         uint64 `json:"id"`
	NickName   string `json:"nick_name"`
	Sex        *int64 `json:"sex"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	InviteCode string `json:"invite_code"`
}

type UserShowReq struct {
	NickName   string `json:"nick_name"`
	Sex        *int64 `json:"sex"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Status     int64  `json:"status"`
	InviteCode string `json:"invite_code"`
}

type UserShowResp struct {
	Id         uint64 `json:"id"`
	NickName   string `json:"nick_name"`
	Sex        *int64 `json:"sex"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Status     int64  `json:"status"`
	InviteCode string `json:"invite_code"`
}
