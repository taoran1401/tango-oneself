package command

import (
	cron2 "taogin/core/cron"
	"taogin/router"
)

type BaseCommand struct {
}

func NewBaseCommand() *BaseCommand {
	return &BaseCommand{}
}

//定时任务相关命令
func (this *BaseCommand) CronStartCommand() {
	var cron *cron2.Cron
	router.InitRouteCrontab(cron)
	cron.Start()
}
