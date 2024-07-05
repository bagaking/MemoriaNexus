package cache

import (
	"context"
	"github.com/google/uuid"
	"github.com/khgame/memstore/cache"
	"github.com/khicago/irr"
	"time"
)

const (
	LockKeyPrefix        cache.Prefix  = "lock:"
	MaxAcquireRetries    int           = 3
	AcquireRetryInterval time.Duration = 50 * time.Millisecond
)

var (
	ErrFailedToReleaseLock = irr.Error("failed to release lock")
	ErrFailedToAcquireLock = irr.Error("failed to acquire lock")
)

// AcquireLock 获取分布式锁，不会自动续期
func AcquireLock(ctx context.Context, key string, expiration time.Duration) (string, error) {
	lockKey := LockKeyPrefix.MakeKey(key)
	lockValue := uuid.New().String()
	for i := 0; i < MaxAcquireRetries; i++ {
		acquired, err := Client().SetNX(ctx, lockKey, lockValue, expiration).Result()
		if err != nil {
			return "", err
		}
		if acquired {
			return lockValue, nil
		}
		time.Sleep(AcquireRetryInterval)
	}
	return "", ErrFailedToAcquireLock
}

// ExecuteWithLock 获取分布式锁并自动续期，执行用户传入的闭包
func ExecuteWithLock(ctx context.Context, key string, renewInterval time.Duration, task func() error) error {
	lockKey := LockKeyPrefix.MakeKey(key)
	lockValue := uuid.New().String()
	for i := 0; i < MaxAcquireRetries; i++ {
		acquired, err := Client().SetNX(ctx, lockKey, lockValue, renewInterval).Result()
		if err != nil {
			return err
		}
		if acquired {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			go startWatchdog(ctx, lockKey, lockValue, renewInterval)
			defer func() {
				_ = ReleaseLock(ctx, key, lockValue)
			}()
			return task()
		}
		time.Sleep(AcquireRetryInterval)
	}
	return ErrFailedToAcquireLock
}

// ReleaseLock 释放分布式锁
func ReleaseLock(ctx context.Context, key, value string) error {
	lockKey := LockKeyPrefix.MakeKey(key)
	script := `
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end`
	result, err := Client().Eval(ctx, script, []string{lockKey}, value).Result()
	if err != nil {
		return irr.Wrap(err, "release lock failed")
	}
	if result.(int64) == 0 {
		return ErrFailedToReleaseLock
	}
	return nil
}

// startWatchdog 启动看门狗协程，定期续期锁
func startWatchdog(ctx context.Context, key, value string, renewInterval time.Duration) {
	exp := renewInterval * 2
	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	for i := 0; i < 1024; i++ { // 避免 while true 防止泄露
		select {
		case <-ticker.C:
			script := `
            if redis.call("GET", KEYS[1]) == ARGV[1] then
                return redis.call("PEXPIRE", KEYS[1], ARGV[2])
            else
                return 0
            end`
			result, err := Client().Eval(ctx, script, []string{key}, value, int64(exp/time.Millisecond)).Result()
			if err != nil || result.(int64) == 0 {
				return // 如果续期失败，停止协程
			}
		case <-ctx.Done():
			return // 如果上下文取消，停止协程
		}
	}
}
