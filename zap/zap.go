package zap

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var (
	log     *zap.Logger
	sugar   *zap.SugaredLogger
	Factory ZapFactory
)

type ZapFactory struct {
	Level       string
	PathName    string
	LogFileName string
	MaxSize     int // 文件大小限制,单位MB
	MaxBackups  int // 最大保留日志文件数量
	MaxAge      int // 日志文件保留天数
}

// Please set the ZapFactory parameter before calling
func InitLog() error {
	if Factory.MaxAge == 0 || Factory.MaxBackups == 0 || Factory.MaxSize == 0 || Factory.LogFileName == "" || Factory.PathName == "" || Factory.Level == "" {
		return errors.New("Insufficient parameters")
	}

	var coreArr []zapcore.Core

	// 获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            // NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info和debug级别,debug级别是最低的
		if Factory.Level == "debug" {
			return lev < zap.ErrorLevel && lev >= zap.DebugLevel
		} else {
			return lev < zap.ErrorLevel && lev >= zap.InfoLevel
		}
	})

	var builder strings.Builder
	builder.WriteString(Factory.PathName)
	builder.WriteString("info_")
	builder.WriteString(Factory.LogFileName)
	infoFileName := builder.String()

	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoFileName,       // 日志文件存放目录，
		MaxSize:    Factory.MaxSize,    // 文件大小限制,单位MB
		MaxBackups: Factory.MaxBackups, // 最大保留日志文件数量
		MaxAge:     Factory.MaxAge,     // 日志文件保留天数
		Compress:   false,              // 是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)

	builder.Reset()
	builder.WriteString(Factory.PathName)
	builder.WriteString("error_")
	builder.WriteString(Factory.LogFileName)
	errorFileName := builder.String()
	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorFileName,      // 日志文件存放目录
		MaxSize:    Factory.MaxSize,    // 文件大小限制,单位MB
		MaxBackups: Factory.MaxBackups, // 最大保留日志文件数量
		MaxAge:     Factory.MaxAge,     // 日志文件保留天数
		Compress:   false,              // 是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	log = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddCallerSkip(1)) // zap.AddCaller()为显示文件名和行号，可省略
	sugar = log.Sugar()

	return nil
}

// 格式化日志
func Infof(s string, v ...interface{}) {
	sugar.Infof(s, v...)
}

// 这个函数接受一个键值对（key-value）形式的参数，用于指定日志消息的上下文信息。这种格式适用于需要记录具有结构化数据的日志
func Infow(s string, v ...interface{}) {
	sugar.Infow(s, v...)
}

func Info(v ...interface{}) {
	sugar.Info(v...)
}

func Debugf(s string, v ...interface{}) {
	sugar.Debugf(s, v...)
}

func Debugw(s string, v ...interface{}) {
	sugar.Debugw(s, v...)
}

func Debug(v ...interface{}) {
	sugar.Debug(v...)
}

func Errorf(s string, v ...interface{}) {
	sugar.Errorf(s, v...)
}

func Errorw(s string, v ...interface{}) {
	sugar.Errorw(s, v...)
}

func Error(v ...interface{}) {
	sugar.Error(v...)
}

func Fatalf(s string, v ...interface{}) {
	sugar.Fatalf(s, v...)
}

func Fatalw(s string, v ...interface{}) {
	sugar.Fatalw(s, v...)
}

func Fatal(v ...interface{}) {
	sugar.Error(v...)
}

func Sync() {
	log.Sync()
}
