package ucredis

import (
	"errors"
	"fmt"
	"time"

	"github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
)

type UcRedisSentinelConnPool struct {
	pool *redis.Pool
}

func NewSentinelPool(name string, addrs []string, cons int) *UcRedisSentinelConnPool {
	return NewSentinelPoolWithPassword(name, addrs, cons, "", 0)
}

func NewSentinelPoolWithPassword(name string, addrs []string, cons int, password string, db int) *UcRedisSentinelConnPool {
	sntnl := &sentinel.Sentinel{
		Addrs:      addrs,
		MasterName: name,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   cons,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			// if need password validation, check the password
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					fmt.Printf("auth fail by password: %s, addr: %s, err: %s \n", password, masterAddr, err.Error())
					c.Close()
					return nil, err
				}
			}

			// if DB is explicitly specified, use the specified DB
			if db > 0 {
				if _, err := c.Do("SELECT", db); err != nil {
					fmt.Printf("select db[%d] failed, err: %s", db, err.Error())
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}

	return &UcRedisSentinelConnPool{pool: redisPool}
}

func (p *UcRedisSentinelConnPool) OpenClient() redis.Conn {
	return p.pool.Get()
}

func (p *UcRedisSentinelConnPool) CloseClient(c redis.Conn) {
	c.Close()
}

func (p *UcRedisSentinelConnPool) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := p.OpenClient()
	if conn == nil {
		return nil, fmt.Errorf("could not get sentinel client from redis connection pool")
	}
	defer p.CloseClient(conn)
	return conn.Do(cmd, args...)
}
