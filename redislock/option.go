package redislock

import (
	"time"
)

/*
const (
	// 默认连接池超过 10 s 释放连接
	DefaultIdleTimeoutSeconds = 10
	// 默认最大激活连接数
	DefaultMaxActive = 100
	// 默认最大空闲连接数
	DefaultMaxIdle = 20

	// 默认的分布式锁过期时间
	DefaultLockExpireSeconds = 3

	// 红锁中每个节点默认的处理超时时间为 50 ms
	DefaultSingleLockTimeout = 50 * time.Millisecond
)
*/

type LockOptions struct {
	isBlock              bool
	watchDogWorkStepTime time.Duration // 看门狗工作时间间隙
	blockWaitingSecond   int64         // 阻塞等待加锁时间
	expireTimeSecond     int64         // 锁过期时间
	watchDogMode         bool          // 是否开启看门狗
}

type LockOption func(*LockOptions)

// WithLock 设置阻塞模式
func WithBlock() LockOption {
	return func(o *LockOptions) {
		o.isBlock = true
	}
}

// 设置阻塞最长等待时间
func WithBlockWaitingSeconds(waitingSeconds int64) LockOption {
	return func(o *LockOptions) {
		o.blockWaitingSecond = waitingSeconds
	}
}

// 设置过期时间
func WithExpireSeconds(expireTime int64) LockOption {
	return func(o *LockOptions) {
		o.expireTimeSecond = expireTime
	}
}

func WithWatchDogMode() LockOption {
	return func(o *LockOptions) {
		o.watchDogMode = true
	}
}

// repairLock 修复分布式锁选项
func repairLock(o *LockOptions) {
	if o.isBlock && o.blockWaitingSecond <= 0 {
		// 默认阻塞等待时间上限为 5 秒
		o.blockWaitingSecond = 5
	}

	if o.expireTimeSecond <= 0 {
		// 关闭看门狗，如果没有设置过期时间
		o.watchDogMode = false
	}

	if o.watchDogMode {
		// 用户开启了看门狗
		// 将看门狗时间设置为 过期时间的 1/4
		o.watchDogWorkStepTime = time.Duration(o.expireTimeSecond) * time.Second / 3
	}

}
