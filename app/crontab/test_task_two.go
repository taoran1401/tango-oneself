package crontab

import (
	"fmt"
	"taogin/core/cron"
)

type TestTaskTwo struct {
	cron.JobInterface
}

func NewTestTaskTwo() *TestTaskTwo {
	return &TestTaskTwo{}
}
func (this *TestTaskTwo) Init()       {}
func (this *TestTaskTwo) Destructor() {}

func (this *TestTaskTwo) Run() {
	fmt.Println("cron: test task two ...")
}
