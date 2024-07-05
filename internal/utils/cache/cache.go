package cache

import (
	"github.com/khgame/memstore/cache"
)

var cli *cache.Cache

func Init() {
	cache.Init("localhost:6379")
	cli = cache.NewPrefixedCli("mem_nexus:")
}

// Client - 得到 *cache.Cache 结构的 Client
//
//	Cache struct {
//		*redis.Client
//	}
func Client() *cache.Cache {
	return cli
}
