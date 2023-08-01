package goredis

import (
	"context"
	"github.com/niuniumart/gosdk/martlog"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	Factory RedisFactory
)

// RedisFactory struct
type RedisFactory struct {
	MaxIdleConn int
	MaxConn     int
	IdleTimeout time.Duration
}

// RedisCli redis cli instance
type RedisCli struct {
	RedisPool *redis.Client
}

const (
	DEFAULT_MAX_CONN      = 3500
	DEFAULT_MAX_IDEL_CONN = 3000
	DEFAULT_IDEL_TIMEOUT  = 300 * time.Second // 300  time.Second
)

// CreateRedisCli func create redis cli
func (p *RedisFactory) CreateRedisCli(pwd, url string) (*RedisCli, error) {
	maxIdleConn := p.MaxIdleConn
	if maxIdleConn == 0 {
		maxIdleConn = DEFAULT_MAX_IDEL_CONN
	}
	maxConn := p.MaxConn
	if maxConn == 0 {
		maxConn = DEFAULT_MAX_CONN
	}
	idleTimeout := p.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = DEFAULT_IDEL_TIMEOUT
	}

	redisPool := redis.NewClient(&redis.Options{
		Addr:            url,
		Password:        pwd,
		PoolSize:        maxConn,
		MaxIdleConns:    maxIdleConn,
		ConnMaxIdleTime: idleTimeout,
	})

	if _, err := redisPool.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	var cli = &RedisCli{
		RedisPool: redisPool,
	}
	return cli, nil
}

// IncrBy func incrBy
func (p *RedisCli) IncrBy(key string, count int64) (int64, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	temp, err := conn.IncrBy(context.Background(), key, count).Result()
	if err != nil {
		martlog.Errorf("get IncrBy error:" + err.Error())
		return 0, err
	}
	return temp, nil
}

// Incr func Increase
func (p *RedisCli) Incr(key string) (string, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	temp, err := conn.Incr(context.Background(), key).Result()
	// 第二个为进制，这里为十进制
	ret := strconv.FormatInt(temp, 10)
	if err != nil {
		martlog.Errorf("get Incr error:" + err.Error())
		return "", err
	}
	return ret, nil
}

// Decr func Decrease
func (p *RedisCli) Decr(pool, key string) (string, error) {

	conn := p.RedisPool.Conn()
	defer conn.Close()

	temp, err := conn.Decr(context.Background(), key).Result()
	if err != nil {
		martlog.Errorf("get Decr error:" + err.Error())
		return "", err
	}
	ret := strconv.FormatInt(temp, 10)
	return ret, nil
}

// DescBy func DescBy
func (p *RedisCli) DecrBy(key string, count int64) (int64, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	temp, err := conn.DecrBy(context.Background(), key, count).Result()
	if err != nil {
		martlog.Errorf("get DecrBy error:" + err.Error())
		return 0, err
	}
	return temp, nil
}

// Set key,value,time
func (p *RedisCli) Set(key, value string, expireTime time.Duration) error {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	_, err := conn.Set(context.Background(), key, value, expireTime).Result()
	if err != nil {
		martlog.Errorf("get DecrBy error:" + err.Error())
		return err
	}

	return nil
}

// SetNX a key/value 已存在返回错误信息
func (p *RedisCli) SetNX(key, value string, expireTime time.Duration) error {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	if err := conn.SetNX(context.Background(), key, value, expireTime).Err(); err != nil {
		martlog.Errorf("get DecrBy error:" + err.Error())
		return err
	}

	return nil
}

// Exists check a key
func (p *RedisCli) Exists(key string) bool {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	exists, err := conn.Exists(context.Background(), key).Result()
	if err != nil {
		martlog.Errorf("")
		return false
	}

	if exists == 1 {
		return true
	}

	return false
}

// Get get a key
func (p *RedisCli) Get(key string) (string, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	ret, err := conn.Get(context.Background(), key).Result()
	if err != nil {
		martlog.Errorf("get Get error:" + err.Error())
		return "", err
	}

	return ret, nil
}

// Delete delete a kye
func (p *RedisCli) Delete(key string) error {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	_, err := conn.Del(context.Background(), key).Result()
	if err != nil {
		martlog.Errorf("get Delete error:" + err.Error())
		return err
	}

	return nil
}

// TTL get the key remain expire time(-2: 表示key不存在 / -1表示key没有过期时间)
func (p *RedisCli) TTL(key string) (int, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	ret, err := conn.TTL(context.Background(), key).Result()
	if err != nil {
		martlog.Errorf("get TTL error:" + err.Error())
		return -3, err
	}

	switch ret {
	case -2:
		martlog.Infof("get TTL key is no exist")
		return -2, nil
	case -1:
		martlog.Infof("get TTL key is no expireTime")
		return -1, nil
	default:
		return int(ret.Seconds()), nil
	}
}

// Expire func Expire（True: key存在，设置过期时间成功 / False: key不存在，设置过期时间失败）
func (p *RedisCli) Expire(key string, expire time.Duration) bool {
	conn := p.RedisPool.Conn()
	defer conn.Close()
	ret, _ := conn.Expire(context.Background(), key, expire).Result()
	return ret
}

// RPush Push back
func (p *RedisCli) RPush(key string, value ...string) (int64, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	ret, err := conn.RPush(context.Background(), key, value).Result()
	if err != nil {
		martlog.Errorf("RPush error:" + err.Error())
		return -1, err
	}
	return ret, nil
}

// LPop pop
func (p *RedisCli) LPop(key string) (string, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	ret, err := conn.LPop(context.Background(), key).Result()
	if err != nil {
		martlog.Errorf("LPop error:" + err.Error())
		return "", err
	}
	return ret, nil
}

// LRange LRange LRange
func (p *RedisCli) LRange(key string, begin, end int) ([]string, error) {
	conn := p.RedisPool.Conn()
	defer conn.Close()

	ret, err := conn.LRange(context.Background(), key, int64(begin), int64(end)).Result()
	if err != nil {
		martlog.Errorf("lrange err", err.Error())
		return nil, err
	}
	return ret, nil
}
