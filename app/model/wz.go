package model

import "time"

type Wz struct {
	Id           uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	Title        string     `json:"title"`
	Status       uint32     `json:"status"`
	WhiteUserId  uint64     `json:"white_user_id"`
	BlackUserId  uint64     `json:"black_user_id"`
	LaunchUserId uint64     `json:"launch_user_id"`
	Ending       uint32     `json:"ending"`
	CreatedAt    *time.Time `json:"created_at"` // 创建时间
	UpdatedAt    *time.Time `json:"updated_at"` // 修改时间
	DeletedAt    *time.Time `json:"deleted_at"` // 删除时间
}

func (Wz) TableName() string {
	return "wz"
}
