package model

import "time"

type Users struct {
	Id              uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	NickName        string     `json:"nick_name"`         // 用户昵称
	WeChatCode      string     `json:"we_chat_code"`      // 微信号
	Sex             *int64     `json:"sex"`               // 性别 0男1女
	Salt            string     `json:"salt"`              // 密码盐
	Password        string     `json:"password"`          // 密码
	Phone           string     `json:"phone"`             // 手机
	Email           string     `json:"email"`             // 邮箱
	Avatar          string     `json:"avatar"`            // 头像
	LastLoginTime   *time.Time `json:"last_login_time"`   // 最后登录时间
	LastOnlineTime  *time.Time `json:"last_online_time"`  // 最后在线时间
	CityCode        string     `json:"city_code"`         // 城市编码，新增
	Status          int64      `json:"status"`            // 用户状态 0正常 1禁用
	Longitude       string     `json:"longitude"`         // 经度
	Latitude        string     `json:"latitude"`          // 纬度
	QqId            string     `json:"qq_id"`             // qq登录的唯一id
	WeChatId        string     `json:"we_chat_id"`        // 微信登录的id
	InviteCode      string     `json:"invite_code"`       // 邀请码
	OnlineStatus    *int64     `json:"online_status"`     // 用户在线状态 0：离线 1：在线 2：忙碌
	OtherInviteCode string     `json:"other_invite_code"` // 填入的别人的邀请码
	Remarks         string     `json:"remarks"`           // 备注,如封号等操作
	CreatedAt       *time.Time `json:"created_at"`        // 创建时间
	UpdatedAt       *time.Time `json:"updated_at"`        // 修改时间
	DeletedAt       *time.Time `json:"deleted_at"`        // 删除时间
}

func (Users) TableName() string {
	return "users"
}
