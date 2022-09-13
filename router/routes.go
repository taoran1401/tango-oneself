package router

import (
	"github.com/gin-gonic/gin"
	"taogin/app/controller"
	"taogin/app/controller/api"
	"taogin/app/middleware"
	ws2 "taogin/core/ws"
)

func InitRoute(RouterGroup *gin.RouterGroup) {
	//websocket router
	ws := RouterGroup.Group("/ws")
	{
		ws.GET("/:token", ws2.WsHandler)
	}

	//设置路由组、中间件
	login := RouterGroup.Group("api")
	{
		//login.GET("index", indexController.Index)
		login.POST("send/phone/code", api.NewLoginController().SendPhoneCode) //发送验证码
		login.POST("login", api.NewLoginController().Login)                   //登录
		//apiGroup.GET("register", controller.LoginController{}.Register)
		//apiGroup.GET("register", controller.LoginController{}.Login)
	}

	business := RouterGroup.Group("api").Use(middleware.NewJwtAuthMiddleware().Handle())
	{
		business.GET("index", controller.NewIndexController().Index)
		business.GET("users/:id", api.NewUserController().Show)
		business.PUT("users/:id", api.NewUserController().Update)
		//apiGroup.GET("register", controller.LoginController{}.Login)
	}
}
