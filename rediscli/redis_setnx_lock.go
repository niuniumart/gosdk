package rediscli

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/niuniumart/gosdk/martlog"

	"github.com/gomodule/redigo/redis"
)

const (
	KEY_PREFIX = "redislock::"
)

type DLock struct {
	RedisCli *RedisCli
}

// Create DLock handler
func CreateDLock(cli *RedisCli) *DLock {
	var p DLock
	p.RedisCli = cli
	return &p
}

func (p *DLock) Lock(key, value string, timeout int) error {
	cli := p.RedisCli
	var lock LockService

	lock = RedisSetnxLock{
		Key:           KEY_PREFIX + key,
		Value:         value,
		TimeoutSecond: timeout,
	}
	err := lock.TryLock(cli)
	if err != nil {
		seelog.Errorf("Lock Error :%s", err)
		return err
	}
	return err
}

// bool 为false表示没有删除，为true表示成功删除
func (p *DLock) Unlock(key, value string) (bool, error) {
	fmt.Println("unlock k v ", key, value)
	cli := p.RedisCli
	var lock LockService
	lock = RedisSetnxLock{
		Key:   KEY_PREFIX + key,
		Value: value,
	}
	return lock.Unlock(cli)
}

type RedisSetnxLock struct {
	Key               string
	Value             string
	TimeoutSecond     int
	MaxWaitMills      int
	MillsBetweenTries int
}

func NanoToMills(nanoTime int64) int64 {
	return nanoTime / 1e6
}

func (r RedisSetnxLock) TryLock(redisCli *RedisCli) error {
	conn := redisCli.RedisPool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SET", r.Key,
		r.Value, "EX", r.TimeoutSecond, "NX"))
	if err != nil {
		martlog.Errorf("redis lock error: ", err)
		return err
	}
	return nil
}

//保证原子性（redis是单线程），避免del删除了，其他client获得的lock
var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return -1
end`)

func (r RedisSetnxLock) Unlock(redisCli *RedisCli) (bool, error) {
	conn := redisCli.RedisPool.Get()
	defer conn.Close()

	value, err := redis.Int(delScript.Do(conn, r.Key, r.Value))
	if err != nil {
		martlog.Errorf("redis script error: %s", err)
		return false, err
	}
	// val==1 删除成功， val==0 之前的锁已过期  val==-1 锁不存在/已经被其他拿走。
	if value == 0 {
		return false, err
	} else if value == -1 {
		return false, err
	}
	return true, nil
}

func (r RedisSetnxLock) SetTimeout(redisCli *RedisCli) error {
	return redisCli.Expire(r.Key, r.TimeoutSecond)
}

type LockService interface {
	TryLock(cli *RedisCli) error
	Unlock(cli *RedisCli) (bool, error)
	SetTimeout(cli *RedisCli) error
}
