package tools

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/niuniumart/gosdk/rediscli"
)

// CacheWrapper cache wrapper
type CacheWrapper struct {
	Cli     *rediscli.RedisCli
	Timeout int
}

//Do 缓存函数执行的结果，主要用于db相关的单个返回值的查询方法的结果缓存。
//@result 缓存的function返回值，必须为指针
//@function 返回类型必须为：单个返回值或者单个返回值+error。如果有error且不为nil，则在Do的error里返回
//@args function的参数
func (p *CacheWrapper) Do(result interface{}, function interface{}, args ...interface{}) error {
	key := GetCacheKey(function, args...)
	seelog.Infof("key:%s", key)

	if p.Cli.Exists(key) {
		epInfoBytes, err := p.Cli.Get(key)
		if err == nil {
			seelog.Infof("get from redis success")
			err := json.Unmarshal(epInfoBytes, result)
			if err == nil {
				seelog.Infof("Unmarshal done")
				return nil
			}
		}
	}

	funcRet, err := FuncProxy(function, args...)
	if err != nil {
		return err
	}

	err = SimpleCopyProperties(result, funcRet)
	if err != nil {
		seelog.Errorf("SimpleCopyProperties failed, %s", err.Error())
		return err
	}

	bytes, err := json.Marshal(funcRet)
	if err != nil {
		seelog.Errorf("Marshal failed, %s", err.Error())
		return err
	}

	err = p.Cli.Set(key, string(bytes), p.Timeout)
	if err != nil {
		seelog.Errorf("redis Set failed, %s", err.Error())
	}

	return nil
}

// GetCacheKey 获取缓存key
func GetCacheKey(function interface{}, args ...interface{}) string {
	key := "DbFunc_" + GetFunctionName(function)

	for _, arg := range args {
		//如果参数中包含db.MysqlCli或者tx(事务)，就跳过该参数
		if _, ok := arg.(*gorm.DB); ok {
			continue
		}
		key = fmt.Sprintf("%s_%v", key, arg)
	}

	return key
}
