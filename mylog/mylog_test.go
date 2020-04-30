package mylog

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	InitLog("d:\\work\\Git\\src\\commonlib\\mylog", "test", 7, 512, 0)
	Debug("测试日志", time.Now())

	time.Sleep(time.Second * 5)
}
