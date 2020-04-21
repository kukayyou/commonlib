package ucmysql

import (
	"github.com/astaxie/beego"
	"strings"
	"time"

	"gnetis.com/golang/core/golib/uclog"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitConnectionPool(cons int, datasrc string) bool {
	return InitConnection(cons, int(beego.AppConfig.DefaultInt64("mysql_idle", 10)), datasrc)
}

func InitConnection(maxopen, maxidle int, datasrc string) bool {
	if maxopen <= 0 || datasrc == "" {
		return false
	}

	uclog.Info("begin to init mysql connection pool")
	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		uclog.Critical("mysql error:register mysql driver fail, %s", err.Error())
		return false
	}

	err = orm.RegisterDataBase("default", "mysql", datasrc, maxidle, maxopen)
	if err != nil {
		uclog.Critical("mysql error:register mysql database fail, %s", err.Error())
		return false
	}
	db, _ := orm.GetDB()
	if db != nil {
		db.SetConnMaxLifetime(time.Hour * 7)
	}
	uclog.Info("end to init mysql connection pool")
	return true
}

func InitConnectionPoolWithName(aliasName string, cons int, datasrc string) bool {
	if cons <= 0 || datasrc == "" || aliasName == "" {
		return false
	}

	uclog.Info("begin to init mysql connection pool")
	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		uclog.Critical("mysql error:register mysql driver fail, %s", err.Error())
		return false
	}

	err = orm.RegisterDataBase(aliasName, "mysql", datasrc)
	if err != nil {
		uclog.Critical("mysql error:register mysql database fail, %s", err.Error())
		return false
	}

	orm.SetMaxOpenConns(aliasName, cons)
	orm.SetMaxIdleConns(aliasName, 10)
	db, _ := orm.GetDB(aliasName)
	if db != nil {
		db.SetConnMaxLifetime(time.Hour * 7)
	}
	uclog.Info("end to init mysql connection pool")
	return true
}

type OrmLogger struct {
}

func (this *OrmLogger) Write(p []byte) (n int, err error) {
	// 如果数据库操作失败，日志信息中包含FAIL标识
	line := string(p)
	if strings.Contains(line, "FAIL") {
		uclog.Critical("mysql error:%s", line)
	} else {
		uclog.Debug(line)
	}
	return len(p), nil
}

func OpenDebug() {
	orm.Debug = true
	orm.DebugLog = orm.NewLog(&OrmLogger{})
}
func CloseDebug() {
	orm.Debug = false
	orm.DebugLog = nil
}
