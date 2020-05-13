package mytypeconv

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	//中文半角空格
	nbsp = []byte{0xC2, 0xA0}
)

func TrimSpace(s string) string {
	s = strings.Replace(s, " ", "", -1)
	return strings.Replace(s, string(nbsp), "", -1)
}

// 字符串反转
func Strrev(src string) string {
	s := []byte(src)
	l := len(s)
	d := make([]byte, l)
	for i := 0; i < l; i++ {
		d[l-i-1] = s[i]
	}
	return string(d)
}

// 获取某天的开始和结束时间戳(秒)
func DayTimestamp(timestamp int64) (int64, int64) {
	if timestamp == 0 {
		timestamp = Create_timesecond()
	}
	t := time.Unix(timestamp, 0)
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	return tm1.Unix(), tm2.Unix()
}

// 获取某天的开始和结束时间戳(秒)---UTC
func DayTimestamp1(timestamp int64) (int64, int64) {
	if timestamp == 0 {
		timestamp = Create_timesecond()
	}
	t := time.Unix(timestamp, 0).UTC()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	return tm1.Unix(), tm2.Unix()
}

func Create_timesecond() int64 {
	t := time.Now()
	return t.Unix()
}
func Create_timenanosecond() int64 {
	t := time.Now()
	return t.Unix()*int64(time.Second) + int64(t.Nanosecond())
}
func Create_timestamp() int64 {
	t := time.Now()
	return t.Unix()*1000 + int64(t.Nanosecond())/(1000*1000)
}

//获取下一天零点时间戳:非UTC时间.以本地时区计算
func CreateNextDayZeroTimestamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	tunix := t.AddDate(0, 0, 1).Unix()

	return tunix*1000 + int64(t.Nanosecond())/(1000*1000)
}

func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	return buf
}

func GetMd5String(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func StrInSlice(val string, arr []string) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

func IsInSlice(val uint64, arr []uint64) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

func IsInSliceInt64(val int64, arr []int64) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}
func IsInSliceInt(val int, arr []int) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

func IsInSliceInt8(e int8, arr []int8) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

func IsInSliceInt16(e int16, arr []int16) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

func IsSubSet(subset []uint64, universalset []uint64) bool {
	for _, v := range subset {
		if !IsInSlice(v, universalset) {
			return false
		}
	}
	return true
}

/*
fix bug
lo:0      Link encap:Local Loopback
          inet addr:192.168.83.58
*/
func GetLocalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inter := range inters {
		if inter.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := inter.Addrs()
		if err != nil {
			return ""
		}

		var localIP string
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			if ip4 := ipnet.IP.To4(); ip4 != nil {
				if ip4[0] == 10 || ip4[0] == 172 || ip4[0] == 192 {
					localIP = ip4.String()
					return localIP
				}
			}
		}
	}
	return ""
}

func ConvertString(str string) string {
	buf := make([]byte, 0, 2*len(str))
	buf = append(buf, '"')

	for _, r := range []rune(str) {
		if r < 128 {
			switch r {
			case '"':
				buf = append(buf, '\\', '"')
			case '\\':
				buf = append(buf, '\\', '\\')
			case '\b':
				buf = append(buf, '\\', 'b')
			case '\f':
				buf = append(buf, '\\', 'f')
			case '\n':
				buf = append(buf, '\\', 'n')
			case '\r':
				buf = append(buf, '\\', 'r')
			case '\t':
				buf = append(buf, '\\', 't')
			case '/':
				buf = append(buf, '\\', '/')
			default:
				if r < 32 {
					buf = append(buf, `\u`...)
					s := fmt.Sprintf("%04x", r)
					buf = append(buf, []byte(s)...)
				} else {
					buf = append(buf, byte(r))
				}
			}
		} else if r < 0xFFFF {
			buf = append(buf, `\u`...)
			s := fmt.Sprintf("%04x", r)
			buf = append(buf, []byte(s)...)
		}
	}
	buf = append(buf, '"')

	return string(buf)
}

type QuoteString string

func (us QuoteString) MarshalJSON() ([]byte, error) {
	buf := []byte(ConvertString(string(us)))
	return buf, nil
}

func ToString(v interface{}) string {
	if r, ok := v.(string); ok {
		return r
	}
	//not string should be convert number to string
	switch v.(type) {
	case uint64:
		return strconv.Itoa(int(v.(uint64)))
	case int64:
		return strconv.Itoa(int(v.(int64)))
	case int:
		return strconv.Itoa((v.(int)))
	case int32:
		return strconv.Itoa(int(v.(int32)))
	case uint32:
		return strconv.Itoa(int(v.(uint32)))
	case float64:
		return strconv.Itoa(int(v.(float64)))
	case int8:
		return strconv.Itoa(int(v.(int8)))
	case uint8:
		return strconv.Itoa(int(v.(uint8)))
	case bool:
		if v.(bool) {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}

func IsInteger(v interface{}) bool {
	switch v.(type) {
	case uint64, int64, int, int32, uint32, float64:
		return true
	}
	return false
}

func ToBool(v interface{}, defaultVal bool) bool {
	switch v.(type) {
	case bool:
		return v.(bool)
	case string:
		s := strings.ToLower(v.(string))
		return s == "true"
	default:
		return defaultVal
	}
}

func ToUint64V2(v interface{}) (r uint64) {
	r, err := ToUint64(v)

	if err != nil {
		return 0
	}

	return r
}

func ToUint64(v interface{}) (uint64, error) {
	switch v.(type) {
	case bool:
		if v.(bool) {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseUint(v.(string), 10, 64)
	case uint64:
		return uint64(v.(uint64)), nil
	case int64:
		return uint64(v.(int64)), nil
	case int:
		return uint64(v.(int)), nil
	case int32:
		return uint64(v.(int32)), nil
	case uint32:
		return uint64(v.(uint32)), nil
	case float64:
		return uint64(v.(float64)), nil
	case int8:
		return uint64(v.(int8)), nil
	case uint8:
		return uint64(v.(uint8)), nil
	default:
		err := fmt.Errorf("cannot convert param to integer")
		return 0, err
	}
}

func ToInt64(v interface{}, defaultVal int64) int64 {
	if v == nil {
		return defaultVal
	}

	switch v.(type) {
	case bool:
		if v.(bool) {
			return 1
		}
		return 0
	case string:
		i, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return defaultVal
		}
		return i
	case uint64:
		return int64(v.(uint64))
	case int64:
		return int64(v.(int64))
	case int:
		return int64(v.(int))
	case int32:
		return int64(v.(int32))
	case uint32:
		return int64(v.(uint32))
	case float64:
		return int64(v.(float64))
	case int8:
		return int64(v.(int8))
	case uint8:
		return int64(v.(uint8))
	}
	return defaultVal
}

func ParseInt(v interface{}, defaultVal int64) int64 {
	return ToInt64(v, defaultVal)
}

func ToInt8(v interface{}, defaultVal int8) int8 {
	return int8(ToInt64(v, int64(defaultVal)))
}

func ParseUint(v interface{}, defaultVal uint64) uint64 {
	if v == nil {
		return defaultVal
	}

	switch v.(type) {
	case bool:
		if v.(bool) {
			return 1
		}
		return 0
	case string:
		i, err := strconv.ParseUint(v.(string), 10, 64)
		if err != nil {
			return defaultVal
		}
		return i
	case uint64:
		return uint64(v.(uint64))
	case int64:
		return uint64(v.(int64))
	case int:
		return uint64(v.(int))
	case int32:
		return uint64(v.(int32))
	case uint32:
		return uint64(v.(uint32))
	case float64:
		return uint64(v.(float64))
	}
	return defaultVal
}

func UniqueIntArray(elements []int64) []int64 {

	encountered := map[int64]bool{}
	result := []int64{}

	for _, v := range elements {
		if encountered[v] != true {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func UniqueUintArray(elements []uint64) []uint64 {

	encountered := map[uint64]bool{}
	result := []uint64{}

	for _, v := range elements {
		if encountered[v] != true {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

// 默认区分大小写
func UniqueStringArray(elements []string, caseInsensitive ...bool) []string {

	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if len(caseInsensitive) > 0 && caseInsensitive[0] {
			v = strings.ToLower(v)
		}
		if encountered[v] != true {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func ToInterfaceArray(intLst []uint64) []interface{} {
	r := make([]interface{}, 0, len(intLst))
	for _, v := range intLst {
		r = append(r, v)
	}
	return r
}

func ToInterfaceArrayInt64(intLst []int64) []interface{} {
	r := make([]interface{}, 0, len(intLst))
	for _, v := range intLst {
		r = append(r, v)
	}
	return r
}

func ToUint64Array(strLst ...string) []uint64 {
	r := make([]uint64, 0, len(strLst))
	for _, v := range strLst {
		r = append(r, ParseUint(v, 0))
	}
	return UniqueUintArray(r)
}
func ToStringArray(intLst ...uint64) []string {
	r := make([]string, 0, len(intLst))
	for _, v := range intLst {
		r = append(r, fmt.Sprintf("%d", v))
	}
	return r
}

// 差集 r = a - b
func DifferenceArray(a []uint64, b []uint64) (r []uint64) {
	encountered := make(map[uint64]bool, len(a))
	for _, v := range a {
		encountered[v] = true
	}
	for _, v := range b {
		if _, found := encountered[v]; found {
			delete(encountered, v)
		}
	}
	for k, _ := range encountered {
		r = append(r, k)
	}
	return
}
func DifferenceIntArray(a []int64, b []int64) (r []int64) {
	encountered := make(map[int64]bool, len(a))
	for _, v := range a {
		encountered[v] = true
	}
	for _, v := range b {
		if _, found := encountered[v]; found {
			delete(encountered, v)
		}
	}
	for k, _ := range encountered {
		r = append(r, k)
	}
	return
}
func DifferenceArrayString(a []string, b []string) (r []string) {
	encountered := make(map[string]bool, len(a))
	for _, v := range a {
		encountered[v] = true
	}
	for _, v := range b {
		if _, found := encountered[v]; found {
			delete(encountered, v)
		}
	}
	for k, _ := range encountered {
		r = append(r, k)
	}
	return
}

func FilterNilData(sList []string) []string {
	r := make([]string, 0, len(sList))
	for _, v := range sList {
		if v != "" {
			r = append(r, v)
		}
	}
	return r
}

func FilterZeroData(iList []uint64) []uint64 {
	r := make([]uint64, 0, len(iList))
	for _, v := range iList {
		if v != 0 {
			r = append(r, v)
		}
	}
	return r
}

func MysqlEscapeString(value string) string {
	var ret []byte
	ret = escapeBytesBackslash([]byte{}, []byte(value))
	return string(ret)

}

// reserveBuffer checks cap(buf) and expand buffer to len(buf) + appendSize.
// If cap(buf) is not enough, reallocate new buffer.
func reserveBuffer(buf []byte, appendSize int) []byte {
	newSize := len(buf) + appendSize
	if cap(buf) < newSize {
		// Grow buffer exponentially
		newBuf := make([]byte, len(buf)*2+appendSize)
		copy(newBuf, buf)
		buf = newBuf
	}
	return buf[:newSize]
}

// escapeBytesBackslash escapes []byte with backslashes (\)
// This escapes the contents of a string (provided as []byte) by adding backslashes before special
// characters, and turning others into specific escape sequences, such as
// turning newlines into \n and null bytes into \0.
func escapeBytesBackslash(buf, v []byte) []byte {
	pos := len(buf)
	buf = reserveBuffer(buf, len(v)*2)

	for _, c := range v {
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = ' '
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos += 1
		}
	}

	return buf[:pos]
}

func JoinUint64Array(list []uint64, sep string) string {
	strList := make([]string, 0, len(list))
	for i := range list {
		strList = append(strList, fmt.Sprintf("%d", list[i]))
	}
	return strings.Join(strList, sep)
}

func JoinInt64Array(list []int64, sep string) string {
	strList := make([]string, 0, len(list))
	for i := range list {
		strList = append(strList, fmt.Sprintf("%d", list[i]))
	}
	return strings.Join(strList, sep)
}

func JoinIntArray(list []int, sep string) string {
	strList := make([]string, 0, len(list))
	for i := range list {
		strList = append(strList, fmt.Sprintf("%d", list[i]))
	}
	return strings.Join(strList, sep)
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// 用于日志输出时对日志进行截断
func TruncateString(value string, length int, ellipsis ...bool) string {
	if length > 0 && len(value) > length {
		value = value[0:length]
		if len(ellipsis) > 0 && ellipsis[0] {
			value += "..."
		}
	}
	return value
}

func Contains(array []string, v string) bool {
	if array != nil {
		for _, p := range array {
			if p == v {
				return true
			}
		}
	}
	return false
}

func ToArray(v interface{}) []byte {
	if value, ok := v.([]byte); ok {
		return value
	} else {
		return nil
	}
}

/*
函数功能：interface转为map
函数入参：
v：待转化数据
函数返回值：转化后的map
*/
func InterfaceToMap(v interface{}) map[string]interface{} {
	re := make(map[string]interface{})
	vt := reflect.ValueOf(v)
	if vt.Kind() == reflect.Map {
		for _, key := range vt.MapKeys() {
			value := vt.MapIndex(key)
			re[ToString(key.Interface())] = value.Interface()
		}
	}
	return re
}
