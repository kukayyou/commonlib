// +build !windows,!nacl,!plan9

package uclog

import (
	"fmt"
	"log/syslog"
	"path"
	"runtime"
	"strings"
)

var (
	w           *syslog.Writer
	processName string
	isSyslog    bool
	logPath     string
)

/*
日志相关配置项：
log_level =debug
log_file  =syslog
options json 字符串
#通过beego file引擎存储访问日志。beego目前支持的引擎有file、console、net、smtp，默认为console，支持同时输出到多个引擎
#文件大小为1GB=1*1024*1024*1024,备份文件名xxx.2013-01-01.2
#file引擎配置参数如下（https://beego.me/docs/module/logs.md）：
#	filename 保存的文件名
#	maxlines 每个文件保存的最大行数，默认值 1000000
#	maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
#	daily 是否按照每天 logrotate，默认是 true
#	maxdays 文件最多保存多少天，默认保存 7 天
#	rotate 是否开启 logrotate，默认是 true
#	level 日志保存的时候的级别，默认是 Trace 级别
#	perm 日志文件权限
log_adapter=file
log_config ={"filename":"/var/log/uclog/nohup.msgserver.out","daily":false,"maxsize":1073741824}

说明
1, log_file为syslog时，
运行日志通过syslog配置，一般为msgserver.log
访问日志通过log_config配置(各个服务器决定是否配置)
2, log_file为目录时
运行日志固定为{log_file}{processName}_runtime.log
访问日志默认为{log_file}{processName}_access.log, 也可以配置log_config进行覆盖
*/
func Initialize(procName string, path string, loglevel string) bool {
	return InitRemoteLog(procName, path, loglevel, "", "")
}

func InitRemoteLog(procName string, path string, loglevel string, network string, raddr string) bool {
	if path == "" || loglevel == "" {
		return false
	}
	if network == "" || raddr == "" {
		network = ""
		raddr = ""
	}
	processName = procName

	if path == "syslog" {
		isSyslog = true
		t, err := syslog.Dial(network, raddr, syslog.LOG_LOCAL0, procName)
		if err != nil {
			fmt.Printf("open common syslog fail, %s", err.Error())
			return false
		}
		w = t

		t, err = syslog.Dial(network, raddr, syslog.LOG_LOCAL6, procName)
		if err != nil {
			fmt.Printf("open monitor syslog fail, %s", err.Error())
			return false
		}
		mw = t

		switch loglevel {
		case "debug":
			level = LOG_LEVEL_DEBUG
		case "info":
			level = LOG_LEVEL_INFO
		case "warn":
			level = LOG_LEVEL_WARN
		case "error":
			level = LOG_LEVEL_ERROR
		case "critical":
			level = LOG_LEVEL_CRITICAL
		default:
			level = LOG_LEVEL_DEBUG
		}

	} else {

		// 类似/var/log/uclog/msgserver.log的形式，提取其路径部分
		if strings.Contains(path, ".") {
			pos := strings.LastIndex(path, "/")
			if pos > 0 {
				path = path[:pos]
			}
		}

		if path[len(path)-1] != '/' {
			path += "/"
		}
		logPath = path

		if !initRuntimeLog(logPath, loglevel) {
			return false
		}
		if !initMonitorLog(logPath) {
			return false
		}
		if !initAccessLog(logPath) {
			return false
		}
	}
	return true
}

func GetProcName() string {
	return processName
}

func getCaller() string {
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		_, filename := path.Split(file)
		fnname := runtime.FuncForPC(pc).Name()
		return fmt.Sprintf("[%s:%d:%s] ", filename, line, fnname)
	}

	return ""
}

func (this *UcLog) genLogPrefix() string {

	this.logIndex++
	// reduce log: del lognum infomation
	headerPrefix := fmt.Sprintf("<requestid:%s>%s ", this.requestId, this.getGid())
	for _, v := range this.header {
		headerPrefix = headerPrefix + "<" + v + "> "
	}

	callerPrefix := ""
	_, file, line, ok := runtime.Caller(2)
	if ok {
		_, filename := path.Split(file)
		// reduce log: del controller.method information
		// fnname := runtime.FuncForPC(pc).Name()
		callerPrefix = fmt.Sprintf("<%s:%d> ", filename, line)
	}

	return headerPrefix + callerPrefix
}

func (this *UcLog) Log_Debug(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_DEBUG {
			w.Debug(this.genLogPrefix() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Debug(this.genLogPrefix()+format, v...)
	}
}

func (this *UcLog) Log_Info(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_INFO {
			w.Info(this.genLogPrefix() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Informational(this.genLogPrefix()+format, v...)
	}
}

func (this *UcLog) Log_Warn(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_WARN {
			w.Warning(this.genLogPrefix() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Warning(this.genLogPrefix()+format, v...)
	}
}

func (this *UcLog) Log_Error(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_ERROR {
			w.Err(this.genLogPrefix() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Error(this.genLogPrefix()+format, v...)
	}
}

func (this *UcLog) Log_Critical(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_CRITICAL {
			w.Crit(this.genLogPrefix() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Critical(this.genLogPrefix()+format, v...)
	}
}

func Debug(format string, v ...interface{}) {

	if isSyslog {
		if w != nil && level >= LOG_LEVEL_DEBUG {
			w.Debug(getGid() + getCaller() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Debug(getGid()+getCaller()+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_INFO {
			w.Info(getGid() + getCaller() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Informational(getGid()+getCaller()+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_WARN {
			w.Warning(getGid() + getCaller() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Warning(getGid()+getCaller()+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_ERROR {
			w.Err(getGid() + getCaller() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Error(getGid()+getCaller()+format, v...)
	}
}

func Critical(format string, v ...interface{}) {
	if isSyslog {
		if w != nil && level >= LOG_LEVEL_CRITICAL {
			w.Crit(getGid() + getCaller() + fmt.Sprintf(format, v...))
		}
	} else {
		log.Critical(getGid()+getCaller()+format, v...)
	}
}
