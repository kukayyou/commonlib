package ucredis

import (
	"gnetis.com/golang/core/golib/uclog"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func initEnv() {
	uclog.Initialize("ucboss", `d:/log/`, "debug")
}

func Test_NewRedisClusterConnPool(t *testing.T) {
	initEnv()

	redisAddrs := []string{
		"10.255.0.75:7000",
		"10.255.0.76:7000",
		"10.255.0.79:7000",
	}

	cluster := NewRedisClusterConnPool(1, redisAddrs)

	testKey := "test:str:set"
	testVal := "testVal"
	_, err := cluster.Do(testKey, "set", testKey, "testVal")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	val, err := redis.String(cluster.Do(testKey, "get", testKey))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if val != testVal {
		t.Errorf("realval: %s not equals to expectVal: %s", val, testVal)
		return
	}

	t.Log("Success")
}

func Test_NewRedisClusterConnPoolWithPassword(t *testing.T) {
	initEnv()

	redisAddrs := []string{
		"192.168.32.241:7000",
		"192.168.32.241:8000",
		"192.168.32.241:9000",
	}

	password := "11111111"
	cluster := NewRedisClusterConnPoolWithPassword(1, redisAddrs, password)

	testKey := "test:str:set"
	testVal := "testVal"
	_, err := cluster.Do(testKey, "set", testKey, "testVal")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	val, err := redis.String(cluster.Do(testKey, "get", testKey))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if val != testVal {
		t.Errorf("realval: %s not equals to expectVal: %s", val, testVal)
		return
	}

	t.Log("Success")
}

// RedisSentinelAddr=192.168.28.221:26379
// RedisSentinelName=my-cluster
// RedisSentinelCons=10

func Test_NewRedisSentinelPool(t *testing.T) {
	initEnv()

	redisAddrs := []string{
		"192.168.35.126:26379",
		//"192.168.28.221:26379",
		"192.168.35.127:26379",
		"192.168.35.128:26379",
	}

	sentinelName := "offEmaster" //"my-cluster"
	cluster := NewSentinelPoolWithPassword(sentinelName, redisAddrs, 1, "Admin@1234", 15)

	testKey := "test:str:set"
	testVal := "testVal"
	_, err := cluster.Do("set", testKey, "testVal")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	//wx:jsapi_ticket:tj7bc645acd7a9964b:ww013f4484e55942f8
	val, err := redis.String(cluster.Do("get", testKey))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if val != testVal {
		t.Errorf("realval: %s not equals to expectVal: %s", val, testVal)
		return
	}

	t.Log("Success")
}
