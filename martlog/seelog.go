// Package seelog for martlog
package martlog

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/google/uuid"
	"github.com/niuniumart/gosdk/requestid"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// init data
func init() {
	seelog.RegisterCustomFormatter("ServiceName", createAppNameFormatter)
	logger, err := seelog.LoggerFromConfigAsString(seelogConfig)
	if err != nil {
		log.Fatal("parsing seelog config file err ", err.Error())
	}
	seelog.ReplaceLogger(logger)
}

const (
	ServiceName   = "ServiceName"
	LogLevelError = "ERROR"
	LogLevelInfo  = "INFO"
	LogLevelDebug = "DEBUG"
	LogLevelWarn  = "WARN"
)

// create app name formatter
func createAppNameFormatter(params string) seelog.FormatterFunc {
	return func(message string, level seelog.LogLevel,
		context seelog.LogContextInterface) interface{} {
		serviceName := os.Getenv(ServiceName)
		if serviceName == "" {
			serviceName = "None"
		}
		return serviceName
	}
}

// Errorf implement errorf
func Errorf(format string, params ...interface{}) {
	prefix := getPrefix(LogLevelError)
	seelog.Errorf(prefix+format+"\n", params...)
}

/* Error implement error
func Error(params ...interface{}) {
	prefix := getPrefix(LogLevelError)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
	seelog.Errorf(prefix, newParams)
}*/

// Infof implement infof
func Infof(format string, params ...interface{}) {
	seelog.Infof(getPrefix(LogLevelInfo)+format, params...)
}

/* Info implement info
func Info(params ...interface{}) {
	prefix := getPrefix(LogLevelInfo)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
	seelog.Infof(newParams...)
}*/

// Debugf implement debug
func Debugf(format string, params ...interface{}) {
	prefix := getPrefix(LogLevelDebug)
	seelog.Debugf(prefix+format, params...)
}

/* Debug implement debug
func Debug(params ...interface{}) {
	prefix := getPrefix(LogLevelDebug)
	seelog.Debug(prefix, params)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
}*/

// Warnf implement warn
func Warnf(format string, params ...interface{}) {
	prefix := getPrefix(LogLevelWarn)
	seelog.Warnf(prefix+format, params...)
}

/*// Warn implement warn
func Warn(params ...interface{}) {
	prefix := getPrefix(LogLevelWarn)
	seelog.Warn(prefix, params)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
}*/

// Flush implement flush
func Flush() {
	seelog.Flush()
}

// implement get prefix
func getPrefix(level string) string {
	callerInfo := getCallerName()
	requestID := requestid.Get()
	if requestID == nil {
		requestIDStr := fmt.Sprintf("%+v", uuid.New())
		requestid.Set(requestIDStr)
	}
	prefix := fmt.Sprintf(":::%s:::%v:::%d:::%s:::", level, requestID, requestid.Goid(), callerInfo)
	return prefix
}

// implement get caller name
func getCallerName() string {
	pc, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s.%d %s", filepath.Base(file), line,
		filepath.Base(runtime.FuncForPC(pc).Name()))
}

var seelogConfig string = `
<seelog minlevel="trace">
	<outputs formatid="fmt_info">
         <filter levels="trace,debug,info,warn,error,critical">
			 <rollingfile formatid="fmt_info" type="size" filename="../log/web.log"  maxsize="104857600" maxrolls="10"/>
         </filter>
         <filter levels="error,critical">
			 <rollingfile formatid="fmt_err" type="size" filename="../log/error/web_error.log"  ` +
	`maxsize="10485760" maxrolls="100"/>
         </filter>
	</outputs>
	<formats>
		<format id="fmt_info" format="%Date(2006-01-02 15:04:05.999):::%Msg%n" />
		<format id="fmt_err" format="%Date(2006-01-02 15:04:05.999):::%Msg%n" />
	</formats>
</seelog>`
