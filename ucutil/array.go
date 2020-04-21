package ucutil

import (
	"bytes"
	"strconv"
	"strings"

	"gnetis.com/golang/core/golib/uclog"
)

func Int64SliceToStringSlice(iArr []int64) []string {
	length := len(iArr)
	sArr := make([]string, length, length)
	for i, v := range iArr {
		sArr[i] = strconv.Itoa(int(v))
	}
	return sArr
}

func StringSliceToInt64Slice(iArr []string) []int64 {
	if len(iArr) == 0 {
		return make([]int64, 0)
	}

	length := len(iArr)
	sArr := make([]int64, length, length)

	for i, item := range iArr {
		if v, err := strconv.ParseInt(item, 10, 64); err == nil {
			sArr[i] = v
		} else {
			sArr[i] = 0
		}
	}
	return sArr
}

func StringSliceToInt8Slice(iArr []string) []int8 {
	if len(iArr) == 0 {
		return make([]int8, 0)
	}

	length := len(iArr)
	sArr := make([]int8, length, length)

	for i, item := range iArr {
		if v, err := strconv.ParseInt(item, 10, 64); err == nil {
			sArr[i] = int8(v)
		} else {
			sArr[i] = 0
		}
	}
	return sArr
}

func ToUint64Slice(arr []int64) []uint64 {
	result := make([]uint64, len(arr))
	for i, v := range arr {
		result[i] = uint64(v)
	}
	return result
}

func ToInt64Slice(arr []uint64) []int64 {
	result := make([]int64, len(arr))
	for i, v := range arr {
		result[i] = int64(v)
	}
	return result
}

func ToInt32Slice(arr []uint64) []int32 {
	result := make([]int32, len(arr))
	for i, v := range arr {
		result[i] = int32(v)
	}
	return result
}

func Int32SliceToUint64Slice(arr []int32) []uint64 {
	result := make([]uint64, len(arr))
	for i, v := range arr {
		result[i] = uint64(v)
	}
	return result
}

func Int64SliceToString(iArr []int64) string {
	length := len(iArr)
	const SPLIT = ","
	var buf bytes.Buffer
	for i, v := range iArr {
		buf.WriteString(strconv.FormatInt(v, 10))
		if i < length-1 {
			buf.WriteString(SPLIT)
		}
	}
	return buf.String()
}

func StringSliceToString(iArr []string) string {
	length := len(iArr)
	const SPLIT = ","
	var buf bytes.Buffer
	for i, v := range iArr {
		buf.WriteString("'")
		buf.WriteString(v)
		buf.WriteString("'")
		if i < length-1 {
			buf.WriteString(SPLIT)
		}
	}
	return buf.String()
}

func Uint64SliceToString(iArr []uint64) string {
	length := len(iArr)
	const SPLIT = ","
	var buf bytes.Buffer
	for i, v := range iArr {
		buf.WriteString(strconv.FormatInt(int64(v), 10))
		if i < length-1 {
			buf.WriteString(SPLIT)
		}
	}
	return buf.String()
}

func StringToInt64Slice(str string) []int64 {
	if str == "" {
		return make([]int64, 0)
	}
	arr := strings.Split(str, ",")
	i, length := 0, len(arr)

	slice := make([]int64, length, length)
	for _, e := range arr {
		v, err := strconv.ParseInt(strings.TrimSpace(e), 10, 64)
		if err != nil {
			continue
		}
		slice[i] = v
		i++
	}
	return slice[:i]
}

// 对两个数组取并集
func UnionInt64Slice(a []int64, b []int64) []int64 {
	if len(a) <= 0 {
		return b
	}
	if len(b) <= 0 {
		return a
	}

	n := 0
	result := make([]int64, len(a)+len(b))
	for _, v := range a {
		result[n] = v
		n++
	}
	for _, v := range b {
		if IsInSliceInt64(v, a) {
			continue
		}
		result[n] = v
		n++
	}
	return result[0:n]

}

// 对两个数组取并集
func UnionUint64Slice(a []uint64, b []uint64) []uint64 {
	if len(a) <= 0 {
		return b
	}
	if len(b) <= 0 {
		return a
	}

	n := 0
	result := make([]uint64, len(a)+len(b))
	for _, v := range a {
		result[n] = v
		n++
	}
	for _, v := range b {
		if IsInSlice(v, a) {
			continue
		}
		result[n] = v
		n++
	}
	return result[0:n]

}

// 对两个数组取并集
func UnionInt8Slice(a []int8, b []int8) []int8 {
	if len(a) <= 0 {
		return b
	}
	if len(b) <= 0 {
		return a
	}

	n := 0
	result := make([]int8, len(a)+len(b))
	for _, v := range a {
		result[n] = v
		n++
	}
	for _, v := range b {
		if IsInSliceInt8(v, a) {
			continue
		}
		result[n] = v
		n++
	}
	return result[0:n]

}

func RemoveFromSlice(arr []int64, val int64) []int64 {
	if len(arr) <= 0 {
		return arr
	}
	result := make([]int64, 0, len(arr))
	for _, v := range arr {
		if val != v {
			result = append(result, v)
		}
	}
	return result
}

func RemoveFromUint64Slice(arr []uint64, val uint64) []uint64 {
	if len(arr) <= 0 {
		return arr
	}
	result := make([]uint64, 0, len(arr))
	for _, v := range arr {
		if val != v {
			result = append(result, v)
		}
	}
	return result
}

func RemoveFromStringSlice(arr []string, val string, args ...interface{}) []string {
	if len(arr) <= 0 {
		return arr
	}

	caseSensitive := true
	if len(args) > 0 {
		caseSensitive = ToBool(args[0], true)
	}

	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if caseSensitive {
			if val != v {
				result = append(result, v)
			}
		} else {
			if strings.ToUpper(val) != strings.ToUpper(v) {
				result = append(result, v)
			}
		}
	}

	return result
}

func IsInStringSlice(arr []string, e string, args ...interface{}) bool {
	if len(arr) == 0 {
		return false
	}

	caseSensitive := true
	if len(args) > 0 {
		caseSensitive = ToBool(args[0], true)
	}

	for _, v := range arr {
		if caseSensitive {
			if v == e {
				return true
			}
		} else {
			if strings.ToUpper(v) == strings.ToUpper(e) {
				return true
			}
		}
	}

	return false
}

// 交集 r = a & b
func IntersectArray(a []int64, b []int64) (r []int64) {
	encountered := make(map[int64]bool, len(a))
	for _, v := range a {
		encountered[v] = false
	}
	for _, v := range b {
		if _, found := encountered[v]; found {
			r = append(r, v)
		}
	}
	return
}

// 交集 r = a & b
func IntersectUintArray(a []uint64, b []uint64) (r []uint64) {
	encountered := make(map[uint64]bool, len(a))
	for _, v := range a {
		encountered[v] = false
	}
	for _, v := range b {
		if _, found := encountered[v]; found {
			r = append(r, v)
		}
	}
	return
}

// 分组
func GroupUintArray(elements []uint64, step int) [][]uint64 {
	if step < 1 {
		return [][]uint64{elements}
	}

	length := len(elements)
	num := length / step
	if length%step > 0 {
		num += 1
	}
	elementLists := make([][]uint64, 0)
	defer func() {
		uclog.Info("elements:%d, step:%d, num:%d, elementLists:%d", len(elements), step, num, len(elementLists))
	}()
	for i := 1; i <= num; i++ {
		start := (i - 1) * step
		end := i * step
		var elementList []uint64
		if i < num {
			elementList = elements[start:end]
		} else {
			elementList = elements[start:]
		}
		if len(elementList) > 0 {
			elementLists = append(elementLists, elementList)
		}
	}
	return elementLists
}
