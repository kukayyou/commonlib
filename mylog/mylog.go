package mylog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	sugarLogger *zap.SugaredLogger
	processName string
)

/*
serverName:server名称
logPath：日志文件保存路径
fileMaxAge：日志保留时长
rotationTime：按时 or 分分割文件
*/

func InitLog(logPath, serverName string, logMaxAge, logMaxSize, logMaxBackUps int, logLevel int8) {
	processName = serverName
	writeSyncer := getLogWriter(logPath, logMaxAge, logMaxSize, logMaxBackUps)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(logLevel))

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logPath string, logMaxAge, logMaxSize, logMaxBackUps int) zapcore.WriteSyncer {
	//fileName := logPath + "\\" + getProcName() + ".log"
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./orderserver.log",
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackUps,
		MaxAge:     logMaxAge,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

//调试日志
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	sugarLogger.Debug(msg, logInfo)
}

//一般日志
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	sugarLogger.Info(msg, logInfo)
}

//告警日志
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	sugarLogger.Warn(msg, logInfo)
}

//错误日志
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	sugarLogger.Error(msg, logInfo)
}

//致命错误日志
func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf("requestId:%s, %s", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	sugarLogger.Fatal(msg, logInfo)
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
