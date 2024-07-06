package cache

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/bagaking/goulp/wlog"
	"github.com/khgame/memstore/cache"
	"github.com/khgame/memstore/prefix"
	"github.com/khicago/irr"
)

const (
	luaRelease = `
local key = KEYS[1]
local identity = ARGV[1]
if redis.call("GET", key) == identity then
	return redis.call("DEL", key)
else
	return 0
end
`

	luaPExpire = `
local key = KEYS[1]
local identity = ARGV[1]
local exp = ARGV[2]
if redis.call("GET", key) == identity then
    return redis.call("EXPIRE", key, exp)
else
	return 0
end
`

	LockKeyPrefix        prefix.Prefix = "lock:"
	MaxAcquireRetries    int           = 5                      // 增加重试次数
	AcquireRetryInterval time.Duration = 100 * time.Millisecond // 增加重试间隔
)

type (
	LockerOps struct {
		*cache.Cache
		deleteLockSha string
		renewLockSha  string
	}
)

var (
	ErrFailedToReleaseLock = irr.Error("failed to release lock")
	ErrFailedToAcquireLock = irr.Error("failed to acquire lock")
	defaultLocker          *LockerOps
	oLocker                sync.Once
)

// Locker 懒加载 LockerOps 实例
func Locker(ctx context.Context) *LockerOps {
	oLocker.Do(func() {
		defaultLocker = &LockerOps{Cache: Client()}
		var err error
		if defaultLocker.deleteLockSha, err = defaultLocker.ScriptLoad(ctx, luaRelease).Result(); err != nil {
			wlog.ByCtx(ctx).Errorf("failed to load lua script, err= %v, script= %v", err, luaRelease)
		}
		if defaultLocker.renewLockSha, err = defaultLocker.ScriptLoad(ctx, luaPExpire).Result(); err != nil {
			wlog.ByCtx(ctx).Errorf("failed to load lua script, err= %v, script= %v", err, luaPExpire)
		}

		wlog.ByCtx(ctx).Infof("locker ops created, deleteLockSha= %s, renewLockSha= %s",
			defaultLocker.deleteLockSha,
			defaultLocker.renewLockSha,
		)
	})

	return defaultLocker
}

// Acquire 获取分布式锁，不会自动续期
func (l *LockerOps) Acquire(ctx context.Context, key string, expiration time.Duration) (string, error) {
	lockKey := LockKeyPrefix.MakeKey(key)
	lockValue := uuid.New().String()
	for i := 0; i < MaxAcquireRetries; i++ {
		acquired, err := l.SetNX(ctx, lockKey, lockValue, expiration).Result()
		if err != nil {
			logrus.WithError(err).Warnf("failed to acquire lock on attempt %d", i+1)
			return "", err
		}
		logrus.Debugf("attempt %d to acquire lock: acquired=%v", i+1, acquired)
		if acquired {
			return lockValue, nil
		}
		time.Sleep(AcquireRetryInterval)
	}
	logrus.Errorf("failed to acquire lock after %d attempts", MaxAcquireRetries)
	return "", ErrFailedToAcquireLock
}

// Execute 获取分布式锁并自动续期，执行用户传入的闭包
func (l *LockerOps) Execute(ctx context.Context, key string, renewInterval time.Duration, task func() error) error {
	lockKey := LockKeyPrefix.MakeKey(key)
	lockValue := uuid.New().String()
	for i := 0; i < MaxAcquireRetries; i++ {
		acquired, err := l.SetNX(ctx, lockKey, lockValue, renewInterval).Result()
		if err != nil {
			logrus.WithError(err).Warnf("failed to acquire lock on attempt %d", i+1)
			return err
		}
		if acquired {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			go l.watchdogKey(ctx, lockKey, lockValue, renewInterval)
			defer func() {
				_ = l.Release(ctx, key, lockValue)
			}()
			return task()
		}
		time.Sleep(AcquireRetryInterval)
	}
	logrus.Errorf("failed to acquire lock after %d attempts", MaxAcquireRetries)
	return ErrFailedToAcquireLock
}

// Release 释放分布式锁
func (l *LockerOps) Release(ctx context.Context, key, value string) error {
	lockKey := LockKeyPrefix.MakeKey(key)
	result, err := l.Eval(ctx, luaRelease, []string{lockKey}, value).Result()
	if err != nil {
		return irr.Wrap(err, "release lock failed")
	}
	wlog.ByCtx(ctx, "release").Infof("release lock, key= %v, value= %v, result= %v", key, value, result)
	if delCount, ok := result.(int64); !ok || delCount == 0 {
		return ErrFailedToReleaseLock
	}
	return nil
}

// watchdogKey 启动看门狗协程，定期续期锁
func (l *LockerOps) watchdogKey(ctx context.Context, key, value string, renewInterval time.Duration) {
	exp := renewInterval * 2
	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	for i := 0; i < 1024; i++ { // 避免 while true 防止泄露
		select {
		case <-ticker.C:
			result, err := l.EvalSha(ctx, l.renewLockSha, []string{key}, value, int64(exp/time.Millisecond)).Result()
			if err != nil || result.(int64) == 0 {
				logrus.WithError(err).Error("failed to renew lock")
				return // 如果续期失败，停止协程
			}
		case <-ctx.Done():
			return // 如果上下文取消，停止协程
		}
	}
}
