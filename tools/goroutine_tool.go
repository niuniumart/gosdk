package tools

import (
	"fmt"
	"github.com/niuniumart/gosdk/martlog"
	"runtime"
	"runtime/debug"
)

//GoRoutine  开启协程，已处理好panic。
//@function 无返回值
//@args function的参数
func GoRoutine(function interface{}, args ...interface{}) {
	go func() {
		if err := recover(); err != nil {
			martlog.Errorf("Routine Panic Recover,Error:%s", err)
			//打印调用栈信息
			debug.PrintStack()
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			martlog.Errorf("panic stack info %s\n", stackInfo)
			return
		}

		FuncProxyWithoutReturn(function, args...)
	}()
}
