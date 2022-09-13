package log

import (
	"go.uber.org/zap"
	"taogin/core/zap_driver"
)

//日志处理(废弃， 现在直接使用zap调用日志输出)
type Logger struct {
	Zap *zap.SugaredLogger
}

func NewLogger(logDir string) *Logger {
	return &Logger{
		Zap: zap_driver.Zap(logDir),
	}
}

//记录日志
//msg 内容
//level 等级
func (this *Logger) Record(content interface{}, level string) {
	switch level {
	case "debug":
		this.Zap.Debug(content)
	case "warn":
		this.Zap.Warn(content)
	case "error":
		this.Zap.Error(content)
	default:
		this.Zap.Info(content)
	}
}

func (this *Logger) Info(content interface{}) {
	this.Record(content, "info")
}

func (this *Logger) Error(content interface{}) {
	this.Record(content, "error")
}

func (this *Logger) Warn(content interface{}) {
	this.Record(content, "warn")
}

func (this *Logger) Debug(content interface{}) {
	this.Record(content, "debug")
}
