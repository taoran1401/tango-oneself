package logic

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"strconv"
	"taogin/app/model"
	"taogin/app/types"
	"taogin/config/global"
	"taogin/core/utils"
	"time"
)

type LoginLogic struct {
}

var (
	SignName     = "rablogs"
	TemplateCode = "SMS_133930080"
)

func NewLoginLogic() *LoginLogic {
	return &LoginLogic{}
}

func (this *LoginLogic) CodeLogin(req *types.LoginReq) (*types.LoginResp, error) {
	fmt.Println(req)
	if !this.CheckCode(req) {
		return nil, errors.New("验证码错误")
	}

	//check user
	user := model.Users{}
	err := global.DB["colorful"].Table(user.TableName()).Where("phone = ?", req.Phone).First(&user).Error
	if err != nil && gorm.ErrRecordNotFound != err {
		global.LOG.Error(err.Error())
		return nil, errors.New("数据库错误")
	}

	var userId uint64 = 0
	if gorm.ErrRecordNotFound == err {
		//用户不存在，增加用户
		userId, err = this.AddUser(req)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		userId = user.Id
	}

	if userId == 0 {
		return nil, errors.New("登录失败")
	}

	//生成token
	token, err := this.CreateUserToken(userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//更新token
	err = this.UpdateUserToken(token, userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &types.LoginResp{
		UserId: userId,
		Token:  token,
	}, nil
}

func (this *LoginLogic) PasswordLogin(req *types.LoginReq) (*types.LoginResp, error) {

	//check user
	user := model.Users{}
	err := global.DB["colorful"].Table(user.TableName()).Where("phone = ?", req.Phone).First(&user).Error
	if err != nil && gorm.ErrRecordNotFound != err {
		return nil, errors.New("数据库错误")
	}

	if gorm.ErrRecordNotFound == err {
		//用户不存在
		return nil, errors.New("账户不存在，请使用短信验证码登录")
	}

	if len(user.Password) == 0 {
		//未设置密码
		return nil, errors.New("账户未设置密码，请使用短信验证码登录")
	}

	//验证密码
	if !utils.ComparePassword(user.Password, req.Password, user.Salt) {
		return nil, errors.New("密码不正确")
	}

	//生成token
	userId := user.Id
	token, err := this.CreateUserToken(userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//更新token
	err = this.UpdateUserToken(token, userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &types.LoginResp{
		UserId: userId,
		Token:  token,
	}, nil
}

func (this *LoginLogic) CreateUserToken(userId uint64) (string, error) {
	baseClaims := types.BaseClaims{
		ID: userId,
	}
	token, err := global.JWT.CreateToken(global.JWT.CreateClaims(baseClaims))
	if err != nil {
		return "", errors.New("token生成失败")
	}
	return token, nil
}

//更新用户token
func (this *LoginLogic) UpdateUserToken(token string, userId uint64) error {
	err := global.DB["colorful"].Table(model.UserToken{}.TableName()).Where("user_id = ?", userId).Update(map[string]interface{}{
		"token": token,
	}).Error
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (this *LoginLogic) CheckCode(req *types.LoginReq) bool {
	cacheCode := global.CACHE.Get("code:" + req.Phone)
	if cacheCode != req.Code {
		return false
	}
	return true
}

func (this *LoginLogic) SendPhoneCode(req *types.SendPhoneCodeReq) (err error) {
	//验证码缓存key
	key := "code:" + req.Phone

	//判断验证码是否过期
	oldCode := global.CACHE.Get(key)
	if len(oldCode) > 0 {
		return errors.New("验证码未过期")
	}

	//gen code
	rand.Seed(time.Now().UnixMilli())
	var code = strconv.Itoa(rand.Intn(8999) + 1000)
	ttl := 60 * time.Second
	if res := global.CACHE.Set(key, code, ttl); !res {
		return errors.New("验证码生成失败")
	}
	templateCode := "SMS_133930080"
	templateParam := this.VerificationSmsTemplate(code)

	//send
	err = global.SMS.Send(req.Phone, templateCode, templateParam)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

//短信模板
func (this *LoginLogic) VerificationSmsTemplate(code string) string {
	return "{\"code\": " + code + "}"
}

func (this *LoginLogic) AddUser(req *types.LoginReq) (uint64, error) {
	var (
		users     model.Users
		userToken model.UserToken
	)

	nowTime := time.Now()

	users.Phone = req.Phone
	users.NickName = "用户" + nowTime.Format("2006-01-02 15:04:05") + strconv.Itoa(rand.Intn(8999)+1000)
	users.LastLoginTime = &nowTime

	//开启事务
	tx := global.DB["colorful"].Begin()

	//add users
	err := tx.Table(users.TableName()).Create(&users).Error
	if err != nil {
		global.LOG.Error(err.Error())
		tx.Rollback()
		return 0, errors.New("数据库错误")
	}

	//add user_token
	userToken.UserId = users.Id
	userToken.LastLoginTime = int64(time.Now().Unix())
	err = tx.Table(userToken.TableName()).Create(&userToken).Error
	if err != nil {
		global.LOG.Error(err.Error())
		tx.Rollback()
		return 0, errors.New("数据库错误")
	}

	tx.Commit()
	return users.Id, nil
}
