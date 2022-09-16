package logic

import (
	"errors"
	"taogin/app/model"
	"taogin/app/types"
	"taogin/config/global"
)

type UserLogic struct {
}

func NewUserLogic() *UserLogic {
	return &UserLogic{}
}

func (this *UserLogic) Show(id string, resp *types.UserShowResp) (*types.UserShowResp, error) {
	var users model.Users
	//用户验证
	if err := this.CheckUser(id); err != nil {
		return nil, errors.New(err.Error())
	}
	//需要返回的数据
	resp.Id = users.Id
	resp.Phone = users.Phone
	resp.NickName = users.NickName
	resp.Sex = users.Sex
	resp.Avatar = users.Avatar
	resp.Status = users.Status
	resp.InviteCode = users.InviteCode
	return resp, nil
}

func (this *UserLogic) Update(id string, req *types.UserShowReq) (bool, error) {
	var users model.Users
	//用户验证
	if err := this.CheckUser(id); err != nil {
		return false, errors.New(err.Error())
	}
	global.DB["colorful"].Table("users").Model(&users).Where("id", id).Updates(map[string]interface{}{
		"nick_name": req.NickName,
		"sex":       req.Sex,
		"phone":     req.Phone,
		"avatar":    req.Avatar,
	})
	return true, nil
}

func (this *UserLogic) CheckUser(id string) error {
	var users model.Users
	//获取数据
	err := global.DB["colorful"].Table("users").Where("id", id).Where("status", 0).First(&users).Error
	if err != nil {
		return errors.New("用户不存在或已禁用")
	}
	return nil
}
