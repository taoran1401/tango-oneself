package command

import (
	"fmt"
	"go.uber.org/zap"
	"taogin/config/global"
	cache2 "taogin/core/cache"
	"taogin/protobuf/protopb"
	"time"
)

type TestCmd struct {
	Log *zap.SugaredLogger
}

func NewTestCmd(log *zap.SugaredLogger) *TestCmd {
	return &TestCmd{Log: log}
}

func (this *TestCmd) Handle() {
	fmt.Println((uint32(protopb.CmdBase_CmdBaseSurrender) << 16))
	return
	nowDate := time.Now().Format("2006-01-02")

	fmt.Println(nowDate)

	//投递生产
	//producer.NewSimpleProducer().Handle()
	//global.LOG.Info("info")
	//this.Log.Info("111")
	//this.Log.Zap.Error("21312312")
	//this.Log.Debug("333")
	//this.Log.Warn("444")

	fmt.Println("test")
	//发送邮件
	//global.EMAIL.Send("工作邮件", "taoran0796@163.com", load.NewEmailTemplate().VerifyTemplate("123"), "")

	cache := cache2.NewCache(global.REDIS["db1"], "c:")
	cache.Set("key03", "value02", 0)
	//cache.Has("key02")
}
