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
	SugarLogger *zap.SugaredLogger
	processName string
)

type LogInfo struct {
	RequestID string `json:"requestId"`
}

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
	var core zapcore.Core
	switch logLevel {
	case -1:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	case 0:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	case 1:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.WarnLevel)
	case 2:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)
	case 3:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.FatalLevel)
	}

	logger := zap.New(core)
	SugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logPath string, logMaxAge, logMaxSize, logMaxBackUps int) zapcore.WriteSyncer {
	fileName := logPath + "\\" + getProcName() + ".log"
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackUps,
		MaxAge:     logMaxAge,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

//调试日志
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(" %s %s ", GetRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	SugarLogger.Debug(msg, logInfo)
}

//一般日志
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(" %s %s ", GetRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	SugarLogger.Info(msg, logInfo)
}

//告警日志
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(" %s %s ", GetRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	SugarLogger.Warn(msg, logInfo)
}

//错误日志
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(" %s %s ", GetRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	SugarLogger.Error(msg, logInfo)
}

//致命错误日志
/*func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf("%s %s ", getRequestId(), getGid())
	logInfo := fmt.Sprintf(format, v...)
	SugarLogger.Fatal(msg, logInfo)
}*/

func (log *LogInfo) SetRequestId() {
	log.RequestID = createRequestId()
}

func (log *LogInfo) GetRequestId() string {
	return log.RequestID
}

//获取请求id
func createRequestId() string {
	t := time.Now()
	return fmt.Sprintf("<requestId:%s-%s-%d.%d.%d>", getProcName(), getLocalIP(), t.Unix(), t.Nanosecond(), rand.Intn(1000))
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
