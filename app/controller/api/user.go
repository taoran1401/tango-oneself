package api

import (
	"github.com/gin-gonic/gin"
	"taogin/app/logic"
	"taogin/app/types"
	"taogin/config/atom"
	"taogin/core/response"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

//根据id查询
func (this *UserController) Show(ctx *gin.Context) {
	var resp types.UserShowResp
	//获取参数
	id := ctx.Param("id")
	//调用逻辑
	logic.NewUserLogic().Show(id, &resp)
	//响应
	response.Success(ctx, resp)
}

//根据id更新
func (this *UserController) Update(ctx *gin.Context) {
	var req *types.UserShowReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, atom.ERROR_CODE_PARAM, "参数错误")
	}

	id := ctx.Param("id")

	logic.NewUserLogic().Update(id, req)

	response.Success(ctx, atom.SUCCESS)
}
