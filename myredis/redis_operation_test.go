package myredis

import (
	"fmt"
	"github.com/kukayyou/commonlib/mylog"
	"os"
	"testing"
	"time"
)

func initLog() {
	mylog.InitLog("d:/log/", "configserver", 50, 1024, 10, -1)
	fmt.Println("initLog")
}

func initRedis() {
	addr := []string{"10.255.0.75:7000"}
	if pool := NewRedisPool(10, addr); pool == nil {
		fmt.Println("initialize redis connection pool fail")
		time.Sleep(time.Millisecond * 100)
		os.Exit(-1)
	}
}

func init() {
	initLog()
	initRedis()
}

func TestRedis(t *testing.T) {
	fmt.Println(SetValue("test:123", "test", 60))
	fmt.Println(GetValue("test:123"))
}
