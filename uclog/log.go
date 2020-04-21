package uclog

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"gnetis.com/golang/core/golib/stack"
	"runtime"
	"strconv"
	"strings"
)

var (
	log   *logs.BeeLogger
	level int
)

const (
	LOG_LEVEL_DEBUG    = 5
	LOG_LEVEL_INFO     = 4
	LOG_LEVEL_WARN     = 3
	LOG_LEVEL_ERROR    = 2
	LOG_LEVEL_CRITICAL = 1
)

func GetLogLevel() int {
	return level
}

func GetLogLevelDesc() string {
	var loglevel string
	switch level {
	case LOG_LEVEL_DEBUG:
		loglevel = "debug"
	case LOG_LEVEL_INFO:
		loglevel = "info"
	case LOG_LEVEL_WARN:
		loglevel = "warn"
	case LOG_LEVEL_ERROR:
		loglevel = "error"
	case LOG_LEVEL_CRITICAL:
		loglevel = "critical"
	}
	return loglevel
}

func SetLogLevel(loglevel string) {
	var beelevel int
	switch loglevel {
	case "debug":
		beelevel = logs.LevelDebug
		level = LOG_LEVEL_DEBUG
	case "info":
		beelevel = logs.LevelInformational
		level = LOG_LEVEL_INFO
	case "warn":
		beelevel = logs.LevelWarning
		level = LOG_LEVEL_WARN
	case "error":
		beelevel = logs.LevelError
		level = LOG_LEVEL_ERROR
	case "critical":
		beelevel = logs.LevelCritical
		level = LOG_LEVEL_CRITICAL
	default:
		beelevel = logs.LevelDebug
		level = LOG_LEVEL_DEBUG
	}
	if log != nil {
		log.SetLevel(beelevel)
	}
}

type UcLog struct {
	gid       string
	requestId string
	logIndex  int
	header    []string
}

func (this *UcLog) GetLogRequestId() string {
	return this.requestId
}
func (this *UcLog) SetLogRequestId(requestId string) {
	this.requestId = requestId
	this.gid = getGid()
}

func (this *UcLog) AddLogHeader(header string) {
	this.header = append(this.header, header)
}

func initRuntimeLog(path, lvl string) bool {
	if lvl == "" {
		return false
	}

	runtimeFile := fmt.Sprintf("%s%s_runtime.log", path, processName)
	// 规范file模式下的日志大小，由默认的256M调整为1G
	config := fmt.Sprintf(`{"filename":"%s", "daily":true, "maxlines":10000000, "maxsize":1073741824}`, runtimeFile)
	log = logs.NewLogger(10000)
	log.EnableFuncCallDepth(true)
	log.SetLogFuncCallDepth(3)
	log.SetLogger("file", config)
	SetLogLevel(lvl)

	return true
}

func InitAccessLog(logCfg string, adapter ...string) bool {
	var adaptername string
	if len(adapter) > 0 {
		adaptername = adapter[0]
	}
	if adaptername == "" {
		adaptername = "file"
	}
	if len(logCfg) <= 0 {
		return false
	}

	if err := beego.SetLogger(adaptername, logCfg); err != nil {
		fmt.Println("access log config error")
		return false
	}
	if adaptername != "console" {
		beego.BeeLogger.DelLogger("console")
	}
	return true
}

func initAccessLog(path string) bool {
	accessFile := fmt.Sprintf("%s%s_access.log", path, processName)
	config := fmt.Sprintf(`{"filename":"%s", "daily":true, "maxlines":10000000, "maxsize":1073741824}`, accessFile)
	InitAccessLog(config, "file")
	return true
}

func (this *UcLog) getGid() string {
	if this.gid == "" {
		this.gid = getGid()
	}
	return this.gid
}

func getGid() string {
	return fmt.Sprintf("<gid:%d>", stack.GoID())
}

const (
	contextLogKey = "logkey"
)

func NewLogContext(ctx context.Context, log *UcLog) context.Context {
	return context.WithValue(ctx, contextLogKey, log)
}
func ExtractLogger(ctx context.Context) (*UcLog, bool) {
	log, ok := ctx.Value(contextLogKey).(*UcLog)
	return log, ok
}

func GoID() (gid int64) {
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
