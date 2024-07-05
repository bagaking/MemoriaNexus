package cache

import (
	"context"
	"github.com/bagaking/goulp/wlog"
	"strconv"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/khgame/memstore/cache"
)

const (
	DefaultLockExpiration = 10 * time.Second
)

type SetOperations struct {
	*cache.Cache
}

func Set() *SetOperations {
	return &SetOperations{Cache: Client()}
}

// Clear 清除指定 key 的缓存。
// ctx: 上下文环境，用于控制请求的取消或超时。
// key: 要清除缓存的键。
// MaxRetryAttempts: 最大重试次数，用于处理删除操作失败的情况。
// 返回操作的错误信息，如果成功则返回 nil。
func (s *SetOperations) Clear(ctx context.Context, key string, MaxRetryAttempts int) error {
	for i := 0; i < MaxRetryAttempts; i++ {
		err := s.Del(ctx, key).Err()
		if err == nil {
			return nil
		}
		time.Sleep(AcquireRetryInterval)
	}
	return nil
}

// GetAll 从 Set 缓存中获取数据
func (s *SetOperations) GetAll(ctx context.Context, key string) ([]string, error) {
	cachedValues, err := s.SMembers(ctx, key).Result() // 缓存优先，先尝试从缓存中获取数据
	if err == nil && len(cachedValues) > 0 {
		return cachedValues, nil
	}

	lockValue, err := AcquireLock(ctx, key, DefaultLockExpiration) // 如果缓存未命中，获取分布式锁，等锁过程中或许已经有写入
	if err != nil {
		return nil, err
	}
	defer ReleaseLock(ctx, key, lockValue)

	// 再次尝试从缓存中获取数据
	return s.SMembers(ctx, key).Result()
}

// Insert 将数据写入 Set 缓存
func (s *SetOperations) Insert(ctx context.Context, key string, Expire time.Duration, values ...string) error {
	lockValue, err := AcquireLock(ctx, key, DefaultLockExpiration)
	if err != nil {
		return err
	}
	defer ReleaseLock(ctx, key, lockValue)

	if err = s.SAdd(ctx, key, values).Err(); err != nil {
		return err
	}
	if Expire < 0 {
		return nil
	}
	return s.Expire(ctx, key, Expire).Err()
}

// GetAllUInt64s 从缓存中获取 uint64 数据
func (s *SetOperations) GetAllUInt64s(ctx context.Context, key string) ([]utils.UInt64, error) {
	cachedValues, err := s.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var values []utils.UInt64
	for _, id := range cachedValues {
		uid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			wlog.ByCtx(ctx, "cache.set.GetAllUInt64s").Errorf("failed to parse uint64 from string, id= %v", id)
			// 处理错误，例如记录日志或跳过此值
			continue
		}
		values = append(values, utils.UInt64(uid))
	}
	return values, nil
}

// InsertUInt64s 将 uint64 数据写入缓存
// 单个元素无法过期，可以考虑 ZSet
func (s *SetOperations) InsertUInt64s(ctx context.Context, key string, Expire time.Duration, values ...utils.UInt64) error {
	var strValues []string
	for _, v := range values {
		strValues = append(strValues, strconv.FormatUint(v.Raw(), 10))
	}
	return s.Insert(ctx, key, Expire, strValues...)
}
