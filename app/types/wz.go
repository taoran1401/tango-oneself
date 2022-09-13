package types

type Wz struct {
	Room   uint64    `json:"Room"`
	Status uint64    `json:"Status"`
	User   *UserBase `json:"user"`
}

type Room struct {
	Id       uint64 `json:"id"`
	TargetId uint64 `json:"target_id"`
	Status   uint64 `json:"status"` //0正常，1已开始
	Cap      uint32 `json:"cap"`    //容量
	Num      uint32 `json:"num"`    //人数
}

type VS struct {
	Room  uint64 `json:"room"`
	Board map[string]bool
}
