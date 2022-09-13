package main

import (
	"runtime"
	"taogin/core"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

const (
	CPU_NUM = 8
)

func main() {
	runtime.GOMAXPROCS(CPU_NUM)
	core.NewServer().Run()
}
