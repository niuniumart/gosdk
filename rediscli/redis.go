// Package rediscli for mgr redis cli
package rediscli

import (
	"errors"
	"fmt"
	"github.com/niuniumart/gosdk/martlog"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	Factory RedisFactory
)
var KeyExistError = errors.New("rediskey exist")

//RedisFactory struct
type RedisFactory struct {
	MaxIdleConn int
	MaxConn     int
	IdleTimeout int
}

//RedisCli redis cli instance
type RedisCli struct {
	RedisPool *redis.Pool
}

const (
	DefaultMaxConn     = 3500
	DefaultMaxIdleConn = 3000
	DefaultIdleTimeout = 300 // use time.Second
)

//CreateRedisCli func create redis cli
func (p *RedisFactory) CreateRedisCli(pwd, url string) (*RedisCli, error) {
	maxIdleConn := p.MaxIdleConn
	if maxIdleConn == 0 {
		maxIdleConn = DefaultMaxIdleConn
	}
	maxConn := p.MaxConn
	if maxConn == 0 {
		maxConn = DefaultMaxConn
	}
	idleTimeout := p.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = DefaultIdleTimeout
	}
	redisPool := &redis.Pool{
		MaxIdle:     maxIdleConn,
		MaxActive:   maxConn,
		IdleTimeout: time.Second * time.Duration(idleTimeout),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}
			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		//应用程序检查健康功能
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}
	var cli = &RedisCli{
		RedisPool: redisPool,
	}
	return cli, nil
}

//IncrBy func incrBy
func (p *RedisCli) IncrBy(key string, count int64) (int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	temp, err := redis.Int64(conn.Do("INCRBY", key, count))
	if err != nil {
		martlog.Errorf("get incr error1:" + err.Error())
		return 0, err
	}
	return temp, nil
}

//Increase func Increase
func (p *RedisCli) Increase(key string) (string, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	temp, err := redis.Int64(conn.Do("INCR", key))

	//第二个为进制，这里为十进制
	ret := strconv.FormatInt(temp, 10)
	if err != nil {
		martlog.Errorf("get incr error:" + err.Error())
		return "", err
	}
	martlog.Infof("key:" + key + "value:" + ret)
	return ret, nil
}

//Decrease func Decrease
func (p *RedisCli) Decrease(pool, key string) (string, error) {

	conn := p.RedisPool.Get()
	defer conn.Close()

	temp, err := redis.Int64(conn.Do("DECR", key))
	ret := strconv.FormatInt(temp, 10)
	if err != nil {
		martlog.Errorf("get incr error1:" + err.Error())
		return "", err
	}
	return ret, nil
}

//DescBy func DescBy
func (p *RedisCli) DescBy(key string, count int64) (int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	temp, err := redis.Int64(conn.Do("DECRBY", key, count))
	if err != nil {
		martlog.Errorf("get incr error1:" + err.Error())
		return 0, err
	}
	return temp, nil
}

//Set key,value,time
func (p *RedisCli) Set(key, value string, time int) error {
	conn := p.RedisPool.Get()
	defer conn.Close()
	if time <= 0 {
		_, err := conn.Do("SET", key, value)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SETEX", key, time, value)
		if err != nil {
			return err
		}
	}
	return nil
}

//SetNX a key/value 已存在返回错误信息
func (p *RedisCli) SetNX(key, value string, time int) error {
	conn := p.RedisPool.Get()
	defer conn.Close()
	if time <= 0 {
		isSuccess, err := conn.Do("SETNX", key, value)
		if err != nil {
			return err
		}
		if isSuccess == 0 {
			return KeyExistError
		}
	} else {
		isSuccess, err := redis.Int64(conn.Do("SETNX", key, value))
		if err != nil {
			return err
		}
		//如果已存在该key值返回错误信息
		if isSuccess == 0 {
			return KeyExistError
		}
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

//Exists check a key
func (p *RedisCli) Exists(key string) bool {
	conn := p.RedisPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

//Get get a key
func (p *RedisCli) Get(key string) ([]byte, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

//Delete delete a kye
func (p *RedisCli) Delete(key string) (bool, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

//LikeDeletes batch delete
func (p *RedisCli) LikeDeletes(key string) error {
	conn := p.RedisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = p.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

//TTL get the key remain expire time
func (p *RedisCli) TTL(key string) (int, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	temp, err := redis.Int64(conn.Do("TTL", key))
	strInt64 := strconv.FormatInt(temp, 10)
	if err != nil {
		martlog.Errorf(err.Error())
		return -3, err
	}
	ret, err := strconv.Atoi(strInt64)

	if err != nil {
		martlog.Errorf(err.Error())
		return -4, err
	}
	return ret, nil
}

//Expire func Expire
func (p *RedisCli) Expire(key string, second int) error {
	conn := p.RedisPool.Get()
	defer conn.Close()

	if second > 0 {
		_, err := conn.Do("EXPIRE", key, second)
		if err != nil {
			return err
		}
	}
	return nil
}

// RPush Push back
func (p *RedisCli) RPush(key string, value string) (int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	temp, err := redis.Int64(conn.Do("RPUSH", key, value))
	if err != nil {
		martlog.Errorf("redis error1:" + err.Error())
		return 0, err
	}
	return temp, nil
}

// RPush Push back
func (p *RedisCli) RPushAndSetExpire(key string, value string, second int) (int64, error) {
	temp, err := p.RPush(key, value)
	if err != nil {
		return temp, err
	}
	err = p.Expire(key, second)
	if err != nil {
		return temp, err
	}
	return temp, nil
}

// LPop pop
func (p *RedisCli) LPop(key string) (string, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()

	temp, err := redis.Bytes(conn.Do("LPOP", key))
	if err != nil {
		martlog.Errorf("redis error1:" + err.Error())
		return "", err
	}
	return string(temp), nil
}

// LRange LRange LRange
func (p *RedisCli) LRange(key string, begin, end int) ([]string, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("lrange", key, begin, end))
	if err != nil {
		martlog.Errorf("lrange err", err.Error())
		return nil, err
	}
	resList := make([]string, 0)
	fmt.Println(values)
	for _, v := range values {
		resList = append(resList, string(v.([]byte)))
	}
	return resList, nil
}
