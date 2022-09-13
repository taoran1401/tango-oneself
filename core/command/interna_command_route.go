package command

import (
	"fmt"
)

func InternalRouteCommand(cmd string) bool {
	flag := false
	switch cmd {
	case "list":
		//列出所有命令
		fmt.Println("list")
	case "help":
		fmt.Println("help")
	case "cron:start":
		//定时任务
		var baseCommand BaseCommand
		baseCommand.CronStartCommand()
	default:
		flag = true
	}
	return flag
}
