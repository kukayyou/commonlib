package myredis

import (
	"fmt"
	"github.com/kukayyou/commonlib/mylog"

	"github.com/garyburd/redigo/redis"
)

var (
	UsersPropertys string = "configserver:usrprop:%d"  //用户属性key
	SitesPropertys string = "configserver:siteprop:%d" //站点属性key
	UserSettings   string = "configserver:usrset:%d"   //用户配置
)

/*
函数功能：写入redis数据
函数入参：
【key】：redis键值
【value】：键值对应的值
【expires】：过期时间
函数返回值：error，返回的错误
*/
func SetValue(key string, value string, expires uint64) error {
	if key == "" || value == "" {
		return fmt.Errorf("SetValue invalid param of key: %s or value: %s", key, value)
	}
	_, err := Do(key, "set", key, value)
	if err != nil {
		return fmt.Errorf("SetValue cache key:%s failed: %s", key, err.Error())
	}
	if expires <= 0 {
		return nil
	}
	_, err = Do(key, "expire", key, expires)
	if err != nil {
		return fmt.Errorf("SetValue set expire time of key: %s failed: %s", key, err.Error())
	}
	mylog.Info("redis cache SetValue success: key = %s,value = %s,expires = %d", key, value, expires)
	return nil
}

/*
函数功能：读取redis数据
函数入参：
【key】：redis键值
函数返回值：string类型，键值对应的值；error，返回的错误
*/
func GetValue(key string) (string, error) {
	value, err := redis.String(Do(key, "get", key))
	if err != nil {
		return "", fmt.Errorf("GetValue not found value for key: %s, %s", key, err.Error())
	}
	mylog.Info("GetValue found by key = %s,value = %s", key, value)
	return value, nil
}

/*
函数功能：删除redis数据
函数入参：
【key】：redis键值
函数返回值：error，返回的错误
*/
func Del(key string) error {
	_, err := Do(key, "del", key)
	if err != nil {
		return fmt.Errorf("del key err: %s", err.Error())
	}
	return nil
}

/*
func GetMultiValue(key []string) ([]string, error) {
	c := OpenClient()
	if c == nil {
		return nil, fmt.Errorf("open client error")
	}
	defer func() {
		CloseClient(c)
	}()
	for _, value := range key {
		c.Pipe(value, value, "GET", value)
	}
	rs, err := c.Commit()
	if err != nil {
		return nil, err
	}
	re := []string{}
	for _, reply := range rs {
		temp, _ := redis.String(reply.Reply, nil)
		re = append(re, temp)
	}
	return re, nil
}

func SetMultiValue(data map[string]interface{}, expires uint64) error {
	c := OpenClient()
	if c == nil {
		return fmt.Errorf("open client error")
	}
	defer func() {
		CloseClient(c)
	}()
	for key, value := range data {
		args := redis.Args{}.Add(key).Add(expires).Add(value)
		c.Pipe(key, key, "SETEX", args)
	}
	_, err := c.Commit()
	if err != nil {
		return err
	}
	return nil
}
*/
