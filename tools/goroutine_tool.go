package tools

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/niuniumart/gosdk/seelog"
)

//GoRoutine  开启协程，已处理好panic。
//@function 无返回值
//@args function的参数
func GoRoutine(function interface{}, args ...interface{}) {
	go func() {
		if err := recover(); err != nil {
			seelog.Errorf("Routine Panic Recover,Error:%s", err)
			//打印调用栈信息
			debug.PrintStack()
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			seelog.Errorf("panic stack info %s\n", stackInfo)
			return
		}

		FuncProxyWithoutReturn(function, args...)
	}()
}
