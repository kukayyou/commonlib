package myredis

import (
	"fmt"
	"github.com/kukayyou/commonlib/mylog"
	"time"

	"github.com/garyburd/redigo/redis"
)

type MyRedisClusterConnPool struct {
	redisChan chan *MyRedisCluster
}

func NewRedisClusterConnPool(consCount int, clusterAddrs []string) (pool *MyRedisClusterConnPool) {
	return NewRedisClusterConnPoolWithPassword(consCount, clusterAddrs, "")
}

// 扩展方法，支持Redis连接的密码认证
func NewRedisClusterConnPoolWithPassword(consCount int, clusterAddrs []string, password string) (pool *MyRedisClusterConnPool) {

	redisChan := make(chan *MyRedisCluster, 10000)

	for i := 0; i < consCount; i++ {
		c := NewUcRedis()
		c.SetPassword(password)
		for _, addr := range clusterAddrs {
			c.AddRedisServer(addr)

		}

		if !c.Dial() {
			return
		}
		select {
		case redisChan <- c:
		default:
		}
	}

	pool = &MyRedisClusterConnPool{redisChan: redisChan}
	return
}

func (p *MyRedisClusterConnPool) OpenClient() *MyRedisCluster {
	select {
	case c := <-p.redisChan:
		return c
	case <-time.After(time.Second * 1):
		mylog.Error("redis error:open client 1s timeout")
		return nil
	}
}

func (p *MyRedisClusterConnPool) CloseClient(c *MyRedisCluster) {
	select {
	case p.redisChan <- c:
	default:
		mylog.Warn("redis error:close client blocked")
	}
}

func (p *MyRedisClusterConnPool) Do(key string, cmd string, args ...interface{}) (interface{}, error) {
	redisClient := p.OpenClient()
	if redisClient == nil {
		return nil, fmt.Errorf("could not get redis client from redis connection pool")
	}
	defer p.CloseClient(redisClient)
	return redisClient.Do(key, cmd, args...)
}

func (p *MyRedisClusterConnPool) Eval(key string, script *redis.Script, args ...interface{}) (interface{}, error) {
	redisClient := p.OpenClient()
	if redisClient == nil {
		return nil, fmt.Errorf("could not get redis client from redis connection pool")
	}
	defer p.CloseClient(redisClient)
	return redisClient.Eval(key, script, args...)
}
