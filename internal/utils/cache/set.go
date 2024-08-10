package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/khgame/memstore/prefix"

	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"

	"github.com/bagaking/goulp/wlog"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/khgame/memstore/cache"
)

const (
	// DefaultLockExpiration = 5 * time.Second // for prod
	DefaultLockExpiration = 300 * time.Second // for dev

	prefixSet prefix.Prefix = "set:"
)

type (
	SetOps struct {
		*cache.Cache
	}

	SupportedLstCacheType interface{ string | utils.UInt64 }
)

func SET() *SetOps {
	return &SetOps{Cache: Client()}
}

// Clear 清除指定 key 的缓存。
// ctx: 上下文环境，用于控制请求的取消或超时。
// key: 要清除缓存的键。
// MaxRetryAttempts: 最大重试次数，用于处理删除操作失败的情况。
// 返回操作的错误信息，如果成功则返回 nil。
func (s *SetOps) Clear(ctx context.Context, key string, MaxRetryAttempts int) error {
	for i := 0; i < MaxRetryAttempts; i++ {
		err := s.Del(ctx, key).Err()
		if err == nil {
			return nil
		}
		time.Sleep(AcquireRetryInterval)
	}
	return nil
}

func (s *SetOps) lock(ctx context.Context, key string) (func(), error) {
	lockKey := prefixSet.MakeKey(key)
	lockValue, err := Locker(ctx).Acquire(ctx, lockKey, DefaultLockExpiration)
	if err != nil {
		return nil, err
	}
	return func() {
		if e := Locker(ctx).Release(ctx, lockKey, lockValue); e != nil {
			wlog.ByCtx(ctx, "cache.set.lock").WithError(e).Errorf("failed to release lock, key= %v", lockKey)
		}
	}, nil
}

// GetAll 从 SET 缓存中获取数据
func (s *SetOps) GetAll(ctx context.Context, key string) ([]string, error) {
	cachedValues, err := s.SMembers(ctx, key).Result() // 缓存优先，先尝试从缓存中获取数据
	if err == nil && len(cachedValues) > 0 {
		return cachedValues, nil
	}

	unlock, err := s.lock(ctx, key) // 如果缓存未命中，获取分布式锁，等锁过程中或许已经有写入
	if err != nil {
		return nil, irr.Wrap(err, "get all by key %s, failed", key)
	}
	defer unlock()

	// 再次尝试从缓存中获取数据
	return s.SMembers(ctx, key).Result()
}

// Insert 将数据写入 SET 缓存
func (s *SetOps) Insert(ctx context.Context, key string, Expire time.Duration, values ...string) error {
	unlock, err := s.lock(ctx, key)
	if err != nil {
		return irr.Wrap(err, "insert %s failed", key)
	}
	defer unlock()

	if err = s.SAdd(ctx, key, typer.SliceMap(values, typer.Any[string])...).Err(); err != nil {
		return irr.Wrap(err, "insert into key %s failed, values= %v", key, values)
	}
	if Expire < 0 {
		return nil
	}
	return s.Expire(ctx, key, Expire).Err()
}

// GetAllUInt64s 从缓存中获取 uint64 数据
func (s *SetOps) GetAllUInt64s(ctx context.Context, key string) ([]utils.UInt64, error) {
	cachedValues, err := s.SMembers(ctx, key).Result()
	if err != nil {
		return nil, irr.Wrap(err, "get all uint64 by key %s, failed", key)
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
func (s *SetOps) InsertUInt64s(ctx context.Context, key string, Expire time.Duration, values ...utils.UInt64) error {
	var strValues []string
	for _, v := range values {
		strValues = append(strValues, strconv.FormatUint(v.Raw(), 10))
	}
	return s.Insert(ctx, key, Expire, strValues...)
}
