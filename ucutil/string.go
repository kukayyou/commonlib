package ucutil

import (
	"encoding/json"
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
)

func ToJson(val interface{}) string {
	if val == nil {
		return ""
	}
	s, err := json.Marshal(val)
	if err != nil {
		fmt.Println("json marshal failed: ", err.Error())
		return ""
	}
	return string(s)
}

func JsonToUint64Slice(str string) (result []uint64) {
	if str == "" {
		result = make([]uint64, 0)
		return
	}

	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		uclog.Info("json unmarshal failed: %s, str: %s", err.Error(), str)
		result = make([]uint64, 0)
		return
	}

	return
}

func SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rsLen := len(rs)
	if start < 0 || start >= rsLen {
		panic(fmt.Sprintf("start is invalid: %d", start))
	}
	if length <= 0 {
		panic(fmt.Sprintf("length is invalid:%d", length))
	}
	end := start + length
	if end > rsLen {
		end = rsLen
	}
	return string(rs[start:end])
}
