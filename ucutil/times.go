package ucutil

import (
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
	"strings"
	"time"
)

func Format(format string, t time.Time) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

func StrToLocalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty input string")
	}
	zoneName, offset := time.Now().Zone()

	zoneValue := offset / 3600 * 100
	if zoneValue > 0 {
		value += fmt.Sprintf(" +%04d", zoneValue)
	} else {
		value += fmt.Sprintf(" -%04d", zoneValue)
	}

	if zoneName != "" {
		value += " " + zoneName
	}
	return StrToTime(value)
}

func StrToTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty input string")
	}
	layouts := []string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700 MST",
		"2006/01/02 15:04:05 -0700",
		"2006/01/02 15:04:05",
		"2006-01-02 -0700 MST",
		"2006-01-02 -0700",
		"2006-01-02",
		"2006/01/02 -0700 MST",
		"2006/01/02 -0700",
		"2006/01/02",
		"2006-01-02 15:04:05 -0700 -0700",
		"2006/01/02 15:04:05 -0700 -0700",
		"2006-01-02 -0700 -0700",
		"2006/01/02 -0700 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}

// 将unix时间戳按照"Y-m-d H:i:s"格式转换成本地时间字符串
func FormatUnixTime(timestamp int64, formats ...string) string {
	t := time.Unix(timestamp, 0)
	var format string
	if len(formats) > 0 {
		format = formats[0]
	} else {
		format = "Y-m-d H:i:s"
	}
	return Format(format, t)
}

func FormatUnixTimeShort(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return Format("Y-m-d", t)
}

// 如果value中没有时区数据则将其当成服务器所在时区的时间来处理
func ToUnixTime(value interface{}) int64 {
	s := ToString(value)
	if len(s) == 0 {
		return 0
	}

	zoneName, offset := time.Now().Zone()

	zoneValue := offset / 3600 * 100
	if zoneValue > 0 {
		s += fmt.Sprintf(" +%04d", zoneValue)
	} else {
		s += fmt.Sprintf(" -%04d", zoneValue)
	}

	if zoneName != "" {
		s += " " + zoneName
	}

	t, err := StrToTime(s)
	if err != nil {
		uclog.Warn("parse time %s failed, error:%s", s, err.Error())
		return 0
	}

	return t.Unix()
}

/*
now: 				1517471501 --- 2018/2/1 15:51:41
TruncateFloor(now): 1517414400 --- 2018/2/1 0:0:0
TruncateCeil(now):	1517500800 --- 2018/2/2 0:0:0
*/
func TruncateFloor(t time.Time) time.Time {
	h, m, s := t.Clock()
	n := t.Nanosecond()
	nano := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(n)
	return t.Add(-nano)
}

func TruncateCeil(t time.Time) time.Time {
	h, m, s := t.Clock()
	n := t.Nanosecond()
	nano := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(n)
	return t.Add(-nano + 24*time.Hour)
}

func TimeStampCompare(o1, o2 int64) int {
	return DateCompare(time.Unix(o1, 0).UTC(), time.Unix(o2, 0).UTC())
}
func DateCompare(t1, t2 time.Time) int {
	year1, month1, day1 := t1.Date()
	year2, month2, day2 := t2.Date()
	switch {
	case year1 > year2:
		return 1
	case year1 == year2:
		switch {
		case month1 > month2:
			return 1
		case month1 == month2:
			switch {
			case day1 > day2:
				return 1
			case day1 == day2:
				return 0
			case day1 < day2:
				return -1
			}
		case month1 < month2:
			return -1
		}
	case year1 < year2:
		return -1
	}
	return 0
}

/*
** 说明：每周从周日计数，每个月的一号如果不是周日的话将剩余星期数算作上个月的。计算过程如下：
** 1、获取当前月的一号是星期几
** 2、获得当前月的第一个周日是几号
** 3、获得参数是第几个星期
 */
func MonthWeek(t time.Time) int {
	year, month, day := t.Date()
	hour, minute, sec := t.Clock()
	monthFstDayTime := time.Date(year, month, 1, hour, minute, sec, 0, t.Location()) //一号
	week := monthFstDayTime.Weekday()

	var monthDay int = 1
	if week != time.Sunday { // 算出第一个星期天到底是几号
		monthDay = 1 + int(time.Sunday) - int(week)
	}
	howManyDays := day - monthDay
	return howManyDays / 7
}

/*
**
 */
func Date(year int, month time.Month, monthWeek int, weekDay time.Weekday, hour, minute, sec int, loc *time.Location) time.Time {
	monthFstDayTime := time.Date(year, month, 1, hour, minute, sec, 0, loc) //一号
	week := monthFstDayTime.Weekday()

	var monthDay int = 1
	if week != time.Sunday { // 算出第一个星期天到底是几号
		monthDay = 1 + int(time.Sunday) - int(week)
	}

	monthDay = monthDay + (monthWeek * 7) + int(weekDay)

	return time.Date(year, month, monthDay, hour, minute, sec, 0, loc)
}
func CalIntervalTime(interval string) int64 {
	ntime := time.Now()
	d, _ := time.ParseDuration(interval)
	return ntime.Add(d).Unix()
}

// 获取当天时间的零点
func GetTodayMinTime() int64 {
	year, month, day := time.Now().Date()
	currentDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentDay.Unix()*1000 + int64(currentDay.Nanosecond()/1000000)
}

// 获取当天最大的时间戳
func GetTodayMaxTime() int64 {
	year, month, day := time.Now().Date()
	currentDay := time.Date(year, month, day, 23, 59, 59, 999, time.Local)
	return currentDay.Unix()*1000 + int64(currentDay.Nanosecond()/1000000)
}

// 获取指定时间所在日期的最大时间
func GetDateMaxTime(dateTime int64) int64 {
	if dateTime <= 0 {
		dateTime = Create_timestamp()
	}
	year, month, day := time.Unix(dateTime/1000, 0).Date()
	currentDay := time.Date(year, month, day, 23, 59, 59, 999, time.Local)
	return currentDay.Unix()*1000 + int64(currentDay.Nanosecond()/1000000)
}

// 返回毫秒时间戳
func GetDateTimestamp(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/(1000*1000)
}
