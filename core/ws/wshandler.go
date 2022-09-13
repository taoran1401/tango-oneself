package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"taogin/app/model"
	"taogin/app/types"
	"taogin/config/atom"
	"taogin/config/global"
	"taogin/core/response"
)

func WsHandler(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		global.LOG.Error("token 为空")
		response.Error(ctx, atom.ERROR_CODE_NO_LOGIN, "未登录")
		return
	}
	claims, err := global.JWT.ParseToken(token)
	if err != nil {
		global.LOG.Error("token：" + err.Error())
		response.Error(ctx, atom.ERROR_CODE_NO_LOGIN, err.Error())
		return
	}

	//获取用户
	user := model.Users{}
	err = global.DB["colorful"].Table(user.TableName()).First(&user, claims.BaseClaims.ID).Error
	if err != nil && gorm.ErrRecordNotFound != err {
		global.LOG.Error(err.Error())
		response.Error(ctx, atom.ERROR_CODE_NO_LOGIN, "用户信息异常")
		return
	}

	//存储用户基本信息
	ctx.Set("user_info", &types.UserBase{
		Id:         user.Id,
		NickName:   user.NickName,
		Sex:        user.Sex,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		InviteCode: user.InviteCode,
	})
	WsServer(ctx, ctx.Writer, ctx.Request)
}
