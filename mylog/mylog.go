package mylog

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Logger       *zap.Logger
	ServerName   string
	LogPath      string
	LogMaxAge    int64
	RotationTime int64
	LogLevel     int8
	processName  string
)

/*
serverName:server名称
logPath：日志文件保存路径
fileMaxAge：日志保留时长
rotationTime：按时 or 分分割文件
*/
func init() {
	processName = ServerName
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	//设置打印的日志级别
	/*var logLevelEnable zapcore.LevelEnabler
	switch logLevel {
	case 1:
		logLevelEnable = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.InfoLevel
		})
	case 2:
		logLevelEnable = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.WarnLevel
		})
	case 3:
		logLevelEnable = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})
	case 4:
		logLevelEnable = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.FatalLevel
		})
	default:
		logLevelEnable = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.DebugLevel+1
		})
	}*/

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	logWriter := getWriter()

	// 最后创建具体的Logger
	infoLevelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.InfoLevel
	})
	warnLevelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.InfoLevel
	})
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(logWriter), infoLevelEnable),
		zapcore.NewCore(encoder, zapcore.AddSync(logWriter), warnLevelEnable),
	)

	Logger = zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
}

//调试日志
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	Logger.Debug(msg,
		zap.String("", toString(logInfo)))
}

//一般日志
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	Logger.Info(msg,
		zap.String("", toString(logInfo)))
}

//告警日志
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	Logger.Warn(msg,
		zap.String("", toString(logInfo)))
}

//错误日志
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	Logger.Error(msg,
		zap.String("", toString(logInfo)))
}

//致命错误日志
func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	Logger.Fatal(msg,
		zap.String("", toString(logInfo)))
}

// 生成rotatelogs的Logger
func getWriter() io.Writer {
	hook, err := rotatelogs.New(
		LogPath+"%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(LogPath),
		rotatelogs.WithMaxAge(time.Duration(LogMaxAge)),          // 按配置保存n天内的日志
		rotatelogs.WithRotationTime(time.Duration(RotationTime)), //按配置时间分割一次日志
	)

	if err != nil {
		panic(err)
	}
	return hook
}

//获取本机ip
func getLocalIP() string {
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

//interface{} 转 string
func toString(v interface{}) string {
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

//获取请求id
func getRequestId() string {
	t := time.Now()
	return fmt.Sprintf("%s-%s-%d.%d.%d", getProcName(), getLocalIP(), t.Unix(), t.Nanosecond(), rand.Intn(1000))
}

//获取server进程名
func getProcName() string {
	return processName
}

//获取goroutine id
func goID() (gid int64) {
	defer func() {
		if err := recover(); err != nil {
			gid = 0
		}
	}()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	s := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	gid, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		gid = 0
		return
	}
	return
}

//获取goroutine id
func getGid() string {
	var r string
	defer func() {
		if err := recover(); err != nil {
			r = ""
		}
	}()
	gid := goID()
	r = fmt.Sprintf("<gid:%d>", gid)
	return r
}
