package crontab

import (
	"fmt"
	"taogin/core/cron"
)

type TestTask struct {
	cron.JobInterface
}

func NewTestTask() *TestTask {
	return &TestTask{}
}
func (this *TestTask) Init()       {}
func (this *TestTask) Destructor() {}

func (this *TestTask) Run() {
	fmt.Println("cron: test task...")
}
