package tools

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

//GetFunctionName 例如对于方法"runtime/debug.FreeOSMemory"，sep不填，则输出该字符串；sep为'.'，则输出"FreeOSMemory"
func GetFunctionName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	//用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

//FuncProxyWithoutReturn case：无返回值
func FuncProxyWithoutReturn(rawFunction interface{}, args ...interface{}) {
	//判断传入的为函数
	function := reflect.ValueOf(rawFunction)
	if function.Kind() != reflect.Func {
		seelog.Errorf("rawFunction is not func type")
		return
	}

	//匹配函数参数与传入参数的个数
	numInFunc := function.Type().NumIn()
	if numInFunc != len(args) {
		seelog.Errorf("numInFunc is equal to args[%d]", len(args))
		return
	}

	values := make([]reflect.Value, 0, numInFunc)
	for i := 0; i < numInFunc; i++ {
		values = append(values, reflect.ValueOf(args[i]))
	}

	//Call的内部实现会对参数类型做匹配。
	//如果有问题，会抛出panic。panic会被捕获，可以查看其内容，得到错误信息
	function.Call(values)
	return
}

//FuncProxy case: 一个返回值;  一个返回值+error; error
func FuncProxy(rawFunction interface{}, args ...interface{}) (interface{}, error) {
	//判断传入的为函数
	function := reflect.ValueOf(rawFunction)
	if function.Kind() != reflect.Func {
		seelog.Errorf("rawFunction is not func type")
		return nil, errors.New("rawFunction is not func type")
	}

	//匹配函数参数与传入参数的个数
	numInFunc := function.Type().NumIn()
	if numInFunc != len(args) {
		seelog.Errorf("numInFunc is equal to args[%d]", len(args))
		return nil, errors.New("numInFunc is equal to args")
	}

	values := make([]reflect.Value, 0, numInFunc)
	for i := 0; i < numInFunc; i++ {
		values = append(values, reflect.ValueOf(args[i]))
	}

	numInReturn := function.Type().NumOut()
	//最多2个返回值；如果有2个返回值，最后一个必须是error
	if numInReturn > 2 {
		seelog.Errorf("not support numInReturn[%d]", numInReturn)
		return nil, errors.New("not support numInReturn")
	} else if numInReturn == 2 {
		lastType := function.Type().Out(numInReturn - 1)
		if !TypeIsError(lastType) {
			seelog.Errorf("last output must be error")
			return nil, errors.New("last output must be error")
		}
	}

	var err error
	var ret interface{}
	outputs := function.Call(values)
	if numInReturn == 1 {
		retValue := outputs[0]
		if TypeIsError(retValue.Type()) {
			if retValue.Interface() != nil {
				err = retValue.Interface().(error)
			}
		} else {
			ret = outputs[0].Interface()
		}
	}
	if numInReturn == 2 {
		ret = outputs[0].Interface()
		if errRet := outputs[1].Interface(); errRet != nil {
			err = errRet.(error)
		}
	}

	return ret, err
}

// TypeIsError 检验是否是error类型
func TypeIsError(inType reflect.Type) bool {
	if !inType.AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return false
	}
	if !inType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return false
	}
	return true
}
