package api

import (
	"github.com/gin-gonic/gin"
	"taogin/app/logic"
	"taogin/app/types"
	"taogin/config/atom"
	"taogin/config/global"
	"taogin/core/response"
	"taogin/core/utils"
)

type LoginController struct {
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

//登录
func (this *LoginController) Login(ctx *gin.Context) {
	var (
		req  types.LoginReq
		resp *types.LoginResp
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		global.LOG.Error("login:" + err.Error())
		response.Error(ctx, atom.ERROR_CODE_PARAM, atom.GetMsgByCode(atom.ERROR_CODE_PARAM))
		return
	}

	//手机验证
	if !utils.VerifyPhoneFormat(req.Phone) {
		response.Error(ctx, atom.ERROR_CODE_PARAM, "手机号码格式错误")
		return
	}

	var err error
	if len(req.Code) > 0 {
		//短信登录
		resp, err = logic.NewLoginLogic().CodeLogin(&req)
		if err != nil {
			response.Error(ctx, atom.ERROR_CODE_PARAM, err.Error())
			return
		}
	} else if len(req.Password) > 0 {
		//密码登录
		resp, err = logic.NewLoginLogic().PasswordLogin(&req)
		if err != nil {
			response.Error(ctx, atom.ERROR_CODE_PARAM, err.Error())
			return
		}
	} else {
		response.Error(ctx, atom.ERROR_CODE_PARAM, "参数错误")
		return
	}

	response.Success(ctx, resp)
}

func (this *LoginController) Register(ctx *gin.Context) {
	var (
		req  types.RegisterReq
		resp *types.RegisterResp
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		global.LOG.Error("register:" + err.Error())
		response.Error(ctx, atom.ERROR_CODE_PARAM, atom.GetMsgByCode(atom.ERROR_CODE_PARAM))
		return
	}

	//手机验证
	if !utils.VerifyPhoneFormat(req.Phone) {
		response.Error(ctx, atom.ERROR_CODE_PARAM, "手机号码格式错误")
		return
	}

	resp, err := logic.NewLoginLogic().Register(&req)
	if err != nil {
		response.Error(ctx, atom.ERROR_CODE_PARAM, err.Error())
		return
	}
	response.Success(ctx, resp)
}

//发送验证码
func (this *LoginController) SendPhoneCode(ctx *gin.Context) {
	var req types.SendPhoneCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		global.LOG.Error(err.Error())
		response.Error(ctx, atom.ERROR_CODE_PARAM, atom.GetMsgByCode(atom.ERROR_CODE_PARAM))
		return
	}

	err := logic.NewLoginLogic().SendPhoneCode(&req)
	if err != nil {
		global.LOG.Error(err.Error())
		response.Error(ctx, atom.ERROR_CODE_SMS, atom.GetMsgByCode(atom.ERROR_CODE_SMS))
		return
	}

	response.Success(ctx, struct{}{})
}
