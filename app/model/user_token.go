package model

import (
	"time"
)

type UserToken struct {
	Id            uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	UserId        uint64     `json:"user_id"`         //用户id
	Token         string     `json:"token"`           // 用户登录token
	LastLoginTime int64      `json:"last_login_time"` // 最后一次登录时间戳
	CreatedAt     *time.Time `json:"created_at"`      // 创建时间
	UpdatedAt     *time.Time `json:"updated_at"`      // 修改时间
	DeletedAt     *time.Time `json:"deleted_at"`      // 删除时间
}

func (UserToken) TableName() string {
	return "user_token"
}
