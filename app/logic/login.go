package logic

import (
	"errors"
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
	if !this.CheckCode(req.Phone, req.Code) {
		return nil, errors.New("验证码错误")
	}

	//check user
	user := model.Users{}
	err := global.DB["default"].Table(user.TableName()).Where("phone = ?", req.Phone).First(&user).Error
	if err != nil && gorm.ErrRecordNotFound != err {
		global.LOG.Error(err.Error())
		return nil, errors.New("数据库错误")
	}

	var userId uint64 = 0
	if gorm.ErrRecordNotFound == err {
		//用户不存在，增加用户
		userId, err = this.AddUser(&types.RegisterReq{
			Phone:           req.Phone,
			Password:        "",
			ConfirmPassword: "",
		})
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
	err := global.DB["default"].Table(user.TableName()).Where("phone = ?", req.Phone).First(&user).Error
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
	err := global.DB["default"].Table(model.UserToken{}.TableName()).Where("user_id = ?", userId).Update(map[string]interface{}{
		"token": token,
	}).Error
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (this *LoginLogic) CheckCode(phone string, code string) bool {
	cacheCode := global.CACHE.Get("code:" + phone)
	if cacheCode != code {
		return false
	}
	return true
}

func (this *LoginLogic) SendPhoneCode(req *types.SendPhoneCodeReq) (err error) {

	//验证码缓存key
	key := "code:" + req.Scene + ":" + req.Phone

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

	switch req.Scene {
	case "register":
		templateCode := "SMS_133930080"
		templateParam := this.VerificationSmsTemplate(code)
		//send
		err = global.SMS.Send(req.Phone, templateCode, templateParam)
		break
	case "update_password":
		templateCode := "SMS_133930080"
		templateParam := this.VerificationSmsTemplate(code)
		//send
		err = global.SMS.Send(req.Phone, templateCode, templateParam)
		break
	default:
		errors.New("参数错误")
		break
	}

	if err != nil {
		return err
	}
	return nil
}

//短信模板
func (this *LoginLogic) VerificationSmsTemplate(code string) string {
	return "{\"code\": " + code + "}"
}

func (this *LoginLogic) AddUser(req *types.RegisterReq) (uint64, error) {
	var (
		users     model.Users
		userToken model.UserToken
	)

	nowTime := time.Now()

	users.Phone = req.Phone
	users.NickName = "用户" + strconv.Itoa(rand.Intn(899999)+100000)
	users.LastLoginTime = &nowTime
	//密码(验证码登录时不需要密码，直接登录成功后设置密码)
	if req.Password != "" {
		salt := utils.RandString(4)
		users.Password = utils.CreatePassword(req.Password, salt)
	}

	//开启事务
	tx := global.DB["default"].Begin()

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

func (this *LoginLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	if !this.CheckCode(req.Phone, req.Code) {
		return nil, errors.New("验证码错误")
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("两次密码不一致")
	}

	//check user
	user := model.Users{}
	err := global.DB["default"].Table(user.TableName()).Where("phone = ?", req.Phone).First(&user).Error
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
		return nil, errors.New("该手机号已注册，请直接登录")
	}

	if userId == 0 {
		return nil, errors.New("注册失败")
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

	return &types.RegisterResp{
		UserId: userId,
		Token:  token,
	}, nil
}
