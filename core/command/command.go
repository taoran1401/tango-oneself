package command

import (
	"os"
	"taogin/router"
)

type CommandProvide struct {
}

//手动调用命令
func (this CommandProvide) Call(args []string) {
	this.Command(args, true)
}

//调用命令
//args: 命令参数
//isCall：是否代码中调用
func (this CommandProvide) Command(args []string, isCall bool) {
	if isCall == true {
		//代码中调用
		args = args
	} else {
		//控制台调用
		args = os.Args[1:]
	}

	//执行命令
	if len(args) > 0 {
		this.LoadRoute(args[0])
	} else {
		panic("没有发现该命令")
	}
}

//加载命令路由
func (this *CommandProvide) LoadRoute(cmd string) {
	//内置命令
	if InternalRouteCommand(cmd) {
		//自定义命令
		router.InitRouteCommand(cmd)
	}
}
