package cron

import (
	"fmt"
	"github.com/robfig/cron"
)

var (
	TaskMap map[string]*TaskJob //任务map
)

type Cron struct {
}

func NewCron() *Cron {
	return &Cron{}
}

//任务详细
type TaskJob struct {
	Rule string
	Job  JobInterface
}

//定时任务
func (this *Cron) Start() {
	fmt.Println("定时任务：")
	cronObj := cron.New()
	for _, v := range TaskMap {
		cronObj.AddJob(v.Rule, v.Job)
	}
	cronObj.Start()
	defer cronObj.Stop()
	select {}
}

//添加任务
func (this *Cron) Add(key string, rule string, job JobInterface) {
	if len(TaskMap) == 0 {
		TaskMap = make(map[string]*TaskJob)
	}
	TaskMap[key] = &TaskJob{
		Rule: rule,
		Job:  job,
	}
}
