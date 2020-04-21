package uclog

import (
	"fmt"
	"path"
	"runtime"
)

var processName string

func Initialize(procName string, path string, loglevel string) bool {
	return InitRemoteLog(procName, path, loglevel, "", "")
}

func InitRemoteLog(procName string, logpath string, loglevel string, network string, raddr string) bool {
	processName = procName
	return initRuntimeLog(logpath, loglevel)
}

func (this *UcLog) genLogPrefix() string {

	this.logIndex++
	headerPrefix := fmt.Sprintf("<requestid:%s>%s ", this.requestId, this.getGid())
	for _, v := range this.header {
		headerPrefix = headerPrefix + v + " "
	}

	callerPrefix := ""
	_, file, line, ok := runtime.Caller(2)
	if ok {
		_, filename := path.Split(file)
		// fnname := runtime.FuncForPC(pc).Name()
		callerPrefix = fmt.Sprintf("[%s:%d] ", filename, line)
	}

	return headerPrefix + callerPrefix
}
func GetProcName() string {
	return processName
}
func (this *UcLog) Log_Debug(format string, v ...interface{}) {
	log.Debug(this.genLogPrefix()+format, v...)
}

func (this *UcLog) Log_Info(format string, v ...interface{}) {
	log.Informational(this.genLogPrefix()+format, v...)
}

func (this *UcLog) Log_Warn(format string, v ...interface{}) {
	log.Warning(this.genLogPrefix()+format, v...)
}

func (this *UcLog) Log_Error(format string, v ...interface{}) {
	log.Error(this.genLogPrefix()+format, v...)
}

func (this *UcLog) Log_Critical(format string, v ...interface{}) {
	log.Critical(this.genLogPrefix()+format, v...)
}

func Debug(format string, v ...interface{}) {
	log.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	log.Informational(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warning(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Error(format, v...)
}

func Critical(format string, v ...interface{}) {
	log.Critical(format, v...)
}
