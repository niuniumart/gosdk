package rediscli

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/niuniumart/gosdk/seelog"
)

// SecKill sec kill struct
type SecKill struct {
	RedisCli *RedisCli
}

// func get total key
func (p *SecKill) getTotalKey(key string) string {
	totalKey := fmt.Sprintf("sec_kill_%s_total", key)
	return totalKey
}

// func get crt key
func (p *SecKill) getCrtKey(key string) string {
	crtKey := fmt.Sprintf("sec_kill_%s_crt", key)
	return crtKey
}

// CreateSecKillHandler func create sec kill handler
func CreateSecKillHandler(cli *RedisCli) *SecKill {
	var p SecKill
	p.RedisCli = cli
	return &p
}

// BuildData func mset data
func (p *SecKill) BuildData(key string, total int64) error {
	totalKey := p.getTotalKey(key)
	crtKey := p.getCrtKey(key)
	conn := p.RedisCli.RedisPool.Get()
	_, err := conn.Do("mset", totalKey, total, crtKey, total)
	return err
}

// Aquire func aquire
func (p *SecKill) Aquire(key string, num int64) error {
	var luaScript = redis.NewScript(1, secKillLua)
	key = p.getCrtKey(key)
	seelog.Infof("key %s\n", key)
	seelog.Infof("num %d\n", num)
	retCode, err := redis.Int64(luaScript.Do(p.RedisCli.RedisPool.Get(), key, fmt.Sprintf("%d", num)))
	if err != nil {
		return err
	}
	return SecKillRetCodeToError(retCode)
}

var (
	SecKillErrKeyNotExist = errors.New("key not exist")
	SecKillErrSelledOut   = errors.New("selled out")
	SecKillErrNotEnough   = errors.New("not enough")
)

// SecKillRetCodeToError func get code to err
func SecKillRetCodeToError(retCode int64) error {
	switch retCode {
	case -1:
		return SecKillErrKeyNotExist
	case -2:
		return SecKillErrKeyNotExist
	case -3:
		return SecKillErrSelledOut
	case -4:
		return SecKillErrNotEnough
	}
	return nil
}

var secKillLua = `
--KEYS[1] 对应商品的key
--ARGV[1] 买多少个
			local key=KEYS[1]
            local subNum = tonumber(ARGV[1])
			if subNum == nil then return -1
			end
            local surplusStock = tonumber(redis.call('get',key))
            if (surplusStock == nil) then return  -2
			end
            if (surplusStock<=0) then return -3
            elseif (subNum > surplusStock) then  return -4
            else
                redis.call('incrby', KEYS[1], -subNum)
                return subNum
            end
`
