package redislock

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync/atomic"
	"time"
)

const RedisLockKeyPrefix = "redisLock:"

var ErrLockAcquiredByOthers = errors.New("lock is acquired by others")
var ErrDeleteLockFailure = errors.New("unlock failed and lock is acquired by others")
var ErrDelayExpire = errors.New("delayExpire failed")

func IsRetryableErr(err error) bool {
	return errors.Is(err, ErrLockAcquiredByOthers)
}

// 基于 redis 实现的分布式锁，不可重入，但保证了对称性
type RedisLock struct {
	LockOptions
	key    string
	token  string
	client *Client

	// 看门狗运作标识
	runningDog int32
	// 停止看门狗
	stopDog context.CancelFunc
}

func NewRedisLock(key string, client *Client, opts ...LockOption) *RedisLock {
	id, _ := uuid.NewUUID()
	r := RedisLock{
		key:    key,
		token:  id.String(),
		client: client,
	}

	for _, opt := range opts {
		opt(&r.LockOptions)
	}

	repairLock(&r.LockOptions)
	return &r
}

// Lock 加锁
func (r *RedisLock) Lock(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			return
		}
		// 加锁成功的情况下，会启动看门狗
		// 关于该锁本身是不可重入的，所以不会出现同一把锁下看门狗重复启动的情况
		r.watchDog(ctx)
	}()

	// 不管是不是阻塞模式，都要先获取一次锁
	err = r.tryLock(ctx)
	if err == nil {
		// 加锁成功
		return nil
	}

	// 非阻塞模式加锁失败直接返回错误
	if !r.isBlock {
		return err
	}

	// 判断错误是否可以允许重试，不可允许的类型则直接返回错误
	if !IsRetryableErr(err) {
		return err
	}

	// 基于阻塞模式持续轮询取锁
	err = r.blockingLock(ctx)
	return
}

func (r *RedisLock) tryLock(ctx context.Context) error {
	// 首先查询锁是否属于自己
	result, err := r.client.pool.SetNX(ctx, r.key, r.token, time.Duration(r.expireTimeSecond)*time.Second).Result()
	if err != nil {
		return err
	}

	// 加锁失败，已经有锁
	if !result {
		return ErrLockAcquiredByOthers
	}

	return nil
}

// 启动看门狗
func (r *RedisLock) watchDog(ctx context.Context) {
	// 1. 非看门狗模式，不处理
	if !r.watchDogMode {
		return
	}

	// 2. 确保之前启动的看门狗已经正常回收
	for !atomic.CompareAndSwapInt32(&r.runningDog, 0, 1) {
	}

	// 3. 启动看门狗
	ctx, r.stopDog = context.WithCancel(ctx)
	go func() {
		defer func() {
			atomic.StoreInt32(&r.runningDog, 0)
		}()
		r.runWatchDog(ctx)
	}()
}

// runWatchDog 看门狗运作
func (r *RedisLock) runWatchDog(ctx context.Context) {
	ticker := time.NewTicker(r.watchDogWorkStepTime)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// 看门狗负责在用户未显式解锁时，持续为分布式锁进行续期
		// 通过 lua 脚本，延期之前会确保保证锁仍然属于自己
		_ = r.DelayExpire(ctx, r.expireTimeSecond)
	}
}

// 更新锁的过期时间，基于 lua 脚本实现操作原子性
func (r *RedisLock) DelayExpire(ctx context.Context, expireSeconds int64) error {
	result, err := r.client.pool.Eval(ctx, LuaCheckAndExpireDistributionLock, []string{r.key}, []interface{}{r.token, expireSeconds}).Bool()
	if err != nil {
		return nil
	}

	if !result {
		return ErrDelayExpire
	}

	return nil
}

func (r *RedisLock) blockingLock(ctx context.Context) error {
	// 阻塞模式等锁时间上限
	timeoutCh := time.After(time.Duration(r.blockWaitingSecond) * time.Second)
	// 轮询 ticker，每隔 50 ms 尝试取锁一次
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		// ctx 终止了
		case <-ctx.Done():
			return fmt.Errorf("lock failed, ctx timeout, err: %w", ctx.Err())
			// 阻塞等锁达到上限时间，返回
		case <-timeoutCh:
			return fmt.Errorf("block waiting time out, err: %w", ErrLockAcquiredByOthers)
		// 放行
		default:
		}

		// 尝试取锁
		err := r.tryLock(ctx)
		if err == nil {
			// 加锁成功，返回结果
			return nil
		}

		// 不可重试类型的错误，直接返回
		if !IsRetryableErr(err) {
			return err
		}
	}

	// 不可达
	return nil
}

// Unlock 解锁. 基于 lua 脚本实现操作原子性.
func (r *RedisLock) Unlock(ctx context.Context) error {
	defer func() {
		// 停止看门狗
		if r.stopDog != nil {
			r.stopDog()
		}
	}()

	result, err := r.client.pool.Eval(ctx, LuaCheckAndDeleteDistributionLock, []string{r.key}, r.token).Bool()
	if err != nil {
		return err
	}

	// 释放成功则为true
	if !result {
		// 释放失败，不是自己的锁
		return ErrDeleteLockFailure
	}

	return nil
}

func (r *RedisLock) getLockKey() string {
	return RedisLockKeyPrefix + r.key
}
