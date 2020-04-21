// +build !windows,!nacl,!plan9

package uclog

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"log/syslog"
)

var (
	mw         *syslog.Writer
	monitorLog *logs.BeeLogger
)

func initMonitorLog(path string) bool {
	if path == "" {
		return false
	}

	filename := fmt.Sprintf("%s%s_monitor.log", path, processName)

	monitorLog = logs.NewLogger(1000)
	monitorLog.EnableFuncCallDepth(false)
	config := fmt.Sprintf(`{"filename":"%s", "maxdays":30}`, filename)
	monitorLog.SetLogger("file", config)

	return true
}

type Monitorlog struct {
	Fields map[string]interface{}
}

func NewMonitorlog() *Monitorlog {
	l := &Monitorlog{}
	l.Fields = make(map[string]interface{})
	return l
}

func (this *Monitorlog) AddField(name string, value interface{}) {
	this.Fields[name] = value
}

func (this *Monitorlog) RemoveField(name string) {
	delete(this.Fields, name)
}

func (this *Monitorlog) Done() {

	b, _ := json.Marshal(this.Fields)

	if isSyslog {
		if mw == nil {
			return
		}
		mw.Info(processName + " " + string(b))
	} else {
		monitorLog.Info(processName + " " + string(b))
	}
}

func (this *Monitorlog) DoneNew(v interface{}) {

	b, _ := json.Marshal(v)

	if isSyslog {
		if mw == nil {
			return
		}
		mw.Info(processName + " " + string(b))
	} else {
		monitorLog.Info(processName + " " + string(b))
	}
}
