package mylog

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	InitLog("d:\\work\\Git\\src\\commonlib\\mylog", "test", 7, 512, 5, 0)
	Debug("测试日志 is :%d", time.Now().Unix())

	time.Sleep(time.Second * 5)
}
