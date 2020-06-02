package myredis

import (
	"fmt"
)

var (
	Pool *MyRedisClusterConnPool //redis 连接池
)

func NewRedisPool(consCount int, clusterAddrs []string) (pool *MyRedisClusterConnPool) {
	pool = NewRedisClusterConnPool(consCount, clusterAddrs)
	Pool = pool
	return
}

func Do(key string, cmd string, args ...interface{}) (interface{}, error) {
	if Pool == nil {
		return nil, fmt.Errorf("allocate RedisClusterConnPool failed")
	}
	return Pool.Do(key, cmd, args...)
}

func OpenClient() *MyRedisCluster {
	if Pool == nil {
		return nil
	}
	return Pool.OpenClient()
}

func CloseClient(c *MyRedisCluster) {
	if Pool == nil {
		return
	}
	Pool.CloseClient(c)
}
