package core

import (
	"fmt"
	"github.com/fvbock/endless"
	"os"
	"os/signal"
	"syscall"
	"taogin/config/global"
	cache2 "taogin/core/cache"
	"taogin/core/command"
	"taogin/core/email_utils"
	"taogin/core/jwt"
	"taogin/core/queue"
	"taogin/core/sms"
	"taogin/core/zap_driver"
	"taogin/router"
	"time"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (this *Server) Run() {
	//加载
	this.Load()
	args := os.Args[1:]
	if len(args) > 0 {
		//控制台模式
		this.ConsoleRun()
	} else {
		//服务模式
		this.ServerRun()
	}
}

func (this *Server) ConsoleRun() {
	var cp command.CommandProvide
	cp.Command([]string{}, false)
}

func (this *Server) ServerRun() {
	// load router
	var routerCore Router
	router := routerCore.Routers()
	//设置静态文件目录和模板文件目录
	router.Static("/assets", "storage/assets")
	router.LoadHTMLGlob("templates/**/*")
	// 优雅的重启服务
	// - 不关闭现有连接（正在运行中的程序）
	// - 新的进程启动并替代旧进程
	// - 新的进程接管新的连接
	// - 连接要随时响应用户的请求，当用户仍在请求旧进程时要保持连接，新用户应请求新进程，不可以出现拒绝请求的情况
	address := fmt.Sprintf(":%s", global.CONFIG.App.Port)
	srv := endless.NewServer(address, router)
	srv.ReadHeaderTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxHeaderBytes = 1 << 16 //左移：相当于1*2^16
	global.LOG.Error(srv.ListenAndServe().Error())
}

//注册服务
func (this *Server) Load() {
	//viper
	global.VIPER = Viper()
	//zap: 非常快的、结构化的，分日志级别的Go日志库
	//global.LOG = log.NewLogger(global.CONFIG.Zap.Director)
	global.LOG = zap_driver.Zap(global.CONFIG.Zap.Director)
	//gorm连接数据库
	global.DB = NewGormMysql().ConnList()
	//redis连接
	global.REDIS = NewRedis().ConnList()
	//mongodb连接
	global.MONGODB = NewMongoDB().ConnList()
	//加载jwt
	global.JWT = jwt.NewJwt(global.CONFIG.Jwt.SigningKey, global.CONFIG.Jwt.ExpiresTime, global.CONFIG.Jwt.Issuer)
	//缓存(这里选择redis需要优化)
	global.CACHE = cache2.NewCache(global.REDIS["db1"], global.CONFIG.Cache.Prefix)
	//cron
	//global.CRON = cron2.NewCron()
	//短信
	global.SMS = sms.NewSms(global.CONFIG.Sms.AccessKeyId, global.CONFIG.Sms.AccessKeySecret, global.CONFIG.Sms.SignName)
	//email
	global.EMAIL = email_utils.NewEmail(
		global.CONFIG.Email.Host,
		global.CONFIG.Email.Port,
		global.CONFIG.Email.Username,
		global.CONFIG.Email.Password,
		global.CONFIG.Email.Nickname,
		global.CONFIG.Email.From,
		global.CONFIG.Email.IsSSL,
	)
	//启动队列
	//this.LoadQueue()
}

//加载队列
func (this *Server) LoadQueue() {
	//启动消费者
	consumer := queue.NewConsumer(global.CONFIG.RabbitMQ)
	router.InitRouteConsumer(consumer)
	consumer.Start()

	//生产者
	producer := queue.NewProducer(global.CONFIG.RabbitMQ[0])
	producer.Init()

	time.Sleep(1 * time.Second)

	// 捕获退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//log.Info("正在关闭...")
	//安全关闭
	consumer.SecurityExit()

	//log.Info("结束.")
}
