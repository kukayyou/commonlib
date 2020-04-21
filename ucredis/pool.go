package ucredis

import (
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
	"time"

	"github.com/garyburd/redigo/redis"
)

type UcRedisClusterConnPool struct {
	redisChan chan *UcRedisCluster
}

func NewRedisClusterConnPool(consCount int, clusterAddrs []string) (pool *UcRedisClusterConnPool) {
	return NewRedisClusterConnPoolWithPassword(consCount, clusterAddrs, "")
}

// 扩展方法，支持Redis连接的密码认证
func NewRedisClusterConnPoolWithPassword(consCount int, clusterAddrs []string, password string) (pool *UcRedisClusterConnPool) {

	redisChan := make(chan *UcRedisCluster, 10000)

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

	pool = &UcRedisClusterConnPool{redisChan: redisChan}
	return
}

func (p *UcRedisClusterConnPool) OpenClient() *UcRedisCluster {
	select {
	case c := <-p.redisChan:
		return c
	case <-time.After(time.Second * 1):
		uclog.Critical("redis error:open client 1s timeout")
		return nil
	}
}

func (p *UcRedisClusterConnPool) CloseClient(c *UcRedisCluster) {
	select {
	case p.redisChan <- c:
	default:
		uclog.Warn("redis error:close client blocked")
	}
}

func (p *UcRedisClusterConnPool) Do(key string, cmd string, args ...interface{}) (interface{}, error) {
	redisClient := p.OpenClient()
	if redisClient == nil {
		return nil, fmt.Errorf("could not get redis client from redis connection pool")
	}
	defer p.CloseClient(redisClient)
	return redisClient.Do(key, cmd, args...)
}

func (p *UcRedisClusterConnPool) Eval(key string, script *redis.Script, args ...interface{}) (interface{}, error) {
	redisClient := p.OpenClient()
	if redisClient == nil {
		return nil, fmt.Errorf("could not get redis client from redis connection pool")
	}
	defer p.CloseClient(redisClient)
	return redisClient.Eval(key, script, args...)
}
