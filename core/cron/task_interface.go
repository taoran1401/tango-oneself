package cron

type JobInterface interface {
	Run()        // 运行逻辑函数
	Init()       // 初始化函数
	Destructor() // 析构函数
}
