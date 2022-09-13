package zap_driver

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"taogin/core/utils"
	"time"
)

func Zap(logDir string) *zap.SugaredLogger {
	//检查目录
	ok, _ := utils.DirExists(logDir)
	if !ok {
		//创建
		os.Mkdir(logDir, 0777)
	}
	//创建core,设置基本参数
	coreInfo := CoreInfo(logDir)   //会打印除debug的其他日志
	coreDebug := CoreDebug(logDir) //全量日志
	coreWarn := CoreWarn(logDir)   //会打印err和warn
	coreError := CoreError(logDir) //错误日志
	//合并core
	core := zapcore.NewTee(coreInfo, coreDebug, coreWarn, coreError)
	//zap.AddCaller(): 输出文件名和行号
	logger := zap.New(core, zap.AddCaller(), zap.WithCaller(true)).Sugar()
	return logger
}

func CoreInfo(logDir string) zapcore.Core {
	//输出到控制台
	//writerSyncer := zapcore.AddSync(os.Stdout)
	//将日志写入文件需要使用zap.New()来传递配置
	//encode: zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()): 编码器(如何写入日志),我们将使用开箱即用的NewJSONEncoder()，并使用预先设置的ProductionEncoderConfig()。
	//WriterSyncer: 指定日志将写到哪里去。我们使用zapcore.AddSync()函数并且将打开的文件句柄传进去
	lumberJackLogger := GetLogWriter(logDir, "info")
	writerSyncer := zapcore.AddSync(lumberJackLogger)
	return zapcore.NewCore(zapcore.NewJSONEncoder(CustomEncodeConfig()), writerSyncer, zapcore.InfoLevel)
}

func CoreDebug(logDir string) zapcore.Core {
	lumberJackLogger := GetLogWriter(logDir, "debug")
	writerSyncer := zapcore.AddSync(lumberJackLogger)
	return zapcore.NewCore(zapcore.NewJSONEncoder(CustomEncodeConfig()), writerSyncer, zapcore.DebugLevel)
}

func CoreWarn(logDir string) zapcore.Core {
	lumberJackLogger := GetLogWriter(logDir, "warn")
	writerSyncer := zapcore.AddSync(lumberJackLogger)
	return zapcore.NewCore(zapcore.NewJSONEncoder(CustomEncodeConfig()), writerSyncer, zapcore.WarnLevel)
}

func CoreError(logDir string) zapcore.Core {
	lumberJackLogger := GetLogWriter(logDir, "error")
	writerSyncer := zapcore.AddSync(lumberJackLogger)
	return zapcore.NewCore(zapcore.NewJSONEncoder(CustomEncodeConfig()), writerSyncer, zapcore.ErrorLevel)
}

//日志写入配置
func GetLogWriter(logDir string, level string) *lumberjack.Logger {
	//日志文件名
	logFile := logDir + "/" + level + ".log"
	nowDate := time.Now().Format("2006-01-02")
	logFile = logDir + "/" + level + "_" + nowDate + ".log"
	//切割配置
	return &lumberjack.Logger{
		Filename: logFile, //日志文件位置
		//MaxSize:  10,      //在进行切割之前，日志文件的最大大小（以MB为单位）
		//MaxBackups: 5,       //保留旧文件的最大个数
		MaxAge:   30,    //保留旧文件的最大天数
		Compress: false, //是否压缩/归档旧文件
	}
}

//自定义日志输出格式
func CustomEncodeConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		//EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}
