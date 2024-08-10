package cache

import (
	"github.com/khgame/memstore/cache"
	"github.com/redis/go-redis/v9"
)

var cli *cache.Cache

func Init(host string) {
	cache.Init(host)
	cli = cache.NewPrefixedCli("mem_nexus:")
}

func InitByRedisClient(cli *redis.Client) {
	c := cache.NewClientByRedisCli(cli)
	cli = c.Client
}

// Client - 得到 *cache.Cache 结构的 Client
//
//	Cache struct {
//		*redis.Client
//	}
func Client() *cache.Cache {
	return cli
}
