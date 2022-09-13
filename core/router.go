package core

import (
	"github.com/gin-gonic/gin"
	default_router "taogin/router"
)

type Router struct {
}

//加载路由
func (this *Router) Routers() (router *gin.Engine) {
	r := gin.Default()
	//全局中间件
	//r.Use(middleware.RbacMiddleware())
	this.BaseRouter(r)
	this.InitRouter(r)
	this.WsInitRouter()
	return r
}

//加载路由
func (this *Router) InitRouter(router *gin.Engine) {
	default_router.InitRoute(router.Group(""))
}

//加载ws路由
func (this *Router) WsInitRouter() {
	default_router.WsInitRoute()
}

//测试路由
func (this *Router) BaseRouter(router *gin.Engine) {
	//test router
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
		})
	})
}
