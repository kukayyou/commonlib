package ucutil

import (
	"encoding/base64"
	"testing"
)

func TestFlateCompress(t *testing.T) {

	// 测试解压
	base64Data := "Y2NgZGRgY2BiYGRmYBZkZmBhZmNgZWDkAIoxCMzgYmBnYGDMZd289icXAweMycPAxcHAyCB98CMHSKPIFA4GZkbD7wu5GVgYGBjYXk7a+LS7nYGHgRuoitFmkhZCFVA6BaJK3uLZnI5nm1c8bet5PnP3071Tn+ze/bRr4bNN8182zHq+bwnQNh7GYKUZHAy8QNWdzAz87GwMAgyMXAxCDDDAw2DEDLSBG2g8A4Ms1Jgdu59tbXy+ovv57HVPO3ufzlkR4PxkRxc32O7QtPy8Et20xNzMnEqrgMy8dLfEvHSFYGcdBceizMQcHYWn+9a92Lv+5ezWl7sn6ij4ZiYX5Rfnp5UoRCZ6pGbqKBSnFmWmWYMNKc6sSrVSMDQpqLBmYAAA"
	zlibData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		t.Error(err.Error())
		return
	}

	thriftData := UnCompressFlateData(zlibData)
	t.Logf("thriftData.length: %d", len(thriftData))

	// 测试压缩
	goZlibData := CompressFlateData(thriftData)
	goBase64Data := base64.StdEncoding.EncodeToString(goZlibData)

	if goBase64Data != base64Data {
		t.Errorf("Unexpected not equals between c++-%d :\n%s\n and go-%d :\n%s", len(base64Data), base64Data, len(goBase64Data), goBase64Data)
		return
	}

	t.Logf("compress str test between c++ and go success.")
}

func TestInt32Max(t *testing.T) {
	const INT32_MAX = int8(^uint8(0) >> 1)
	t.Logf("maxValue: %d", INT32_MAX)
}
