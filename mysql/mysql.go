package mysql

import (
	"github.com/kukayyou/commonlib/mylog"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitConnectionPool(cons int, datasrc string) bool {
	return InitConnection(cons, 10, datasrc)
}

func InitConnection(maxopen, maxidle int, datasrc string) bool {
	if maxopen <= 0 || datasrc == "" {
		return false
	}

	mylog.Info("begin to init mysql connection pool")
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		mylog.Error("mysql error:register mysql driver fail, %s", err.Error())
		return false
	}

	err = orm.RegisterDataBase("default", "mysql", datasrc, maxidle, maxopen)
	if err != nil {
		mylog.Error("mysql error:register mysql database fail, %s", err.Error())
		return false
	}
	db, _ := orm.GetDB()
	if db != nil {
		db.SetConnMaxLifetime(time.Hour * 7)
	}
	mylog.Info("end to init mysql connection pool")
	return true
}

func InitConnectionPoolWithName(aliasName string, cons int, datasrc string) bool {
	if cons <= 0 || datasrc == "" || aliasName == "" {
		return false
	}

	mylog.Info("begin to init mysql connection pool")
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		mylog.Error("mysql error:register mysql driver fail, %s", err.Error())
		return false
	}

	err = orm.RegisterDataBase(aliasName, "mysql", datasrc)
	if err != nil {
		mylog.Error("mysql error:register mysql database fail, %s", err.Error())
		return false
	}

	orm.SetMaxOpenConns(aliasName, cons)
	orm.SetMaxIdleConns(aliasName, 10)
	db, _ := orm.GetDB(aliasName)
	if db != nil {
		db.SetConnMaxLifetime(time.Hour * 7)
	}
	mylog.Info("end to init mysql connection pool")
	return true
}

type OrmLogger struct {
}

func (this *OrmLogger) Write(p []byte) (n int, err error) {
	// 如果数据库操作失败，日志信息中包含FAIL标识
	line := string(p)
	if strings.Contains(line, "FAIL") {
		mylog.Error("mysql error:%s", line)
	} else {
		mylog.Debug(line)
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
