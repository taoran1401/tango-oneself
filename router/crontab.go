package router

import (
	"taogin/app/crontab"
	"taogin/config/global"
	"taogin/core/cron"
)

func InitRouteCrontab(cron *cron.Cron) {
	//参数说明： key, 规则, 任务
	cron.Add("TestTask", global.CONFIG.Crontab.TestTask, crontab.NewTestTask())
	cron.Add("TestTaskTwo", global.CONFIG.Crontab.TestTask, crontab.NewTestTaskTwo())
}
