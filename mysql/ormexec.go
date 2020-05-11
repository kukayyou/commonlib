package ucmysql

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kukayyou/commonlib/mytypeconv"
)

func GetOrmer(ormer orm.Ormer) orm.Ormer {
	if ormer == nil {
		return orm.NewOrm()
	} else {
		return ormer
	}
}

func GetString(param orm.Params, key string) string {
	if val, ok := param[key]; ok {
		return mytypeconv.ToString(val)
	} else {
		return ""
	}
}

func ExecForDelete(sql string, o orm.Ormer) (int64, error) {
	return ExecForUpdate(sql, o)
}

func ExecForUpdate(sql string, o orm.Ormer) (int64, error) {
	rs, err := o.Raw(sql).Exec()
	if err != nil {
		return 0, fmt.Errorf("update sql execute error: %s, %s", err.Error(), sql)
	}
	num, err := rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read update count from rs failed: %s", err.Error())
	}
	return num, nil
}

func ExecForInsert(sql string, o orm.Ormer) (int64, error) {
	rs, err := o.Raw(sql).Exec()
	if err != nil {
		return 0, fmt.Errorf("insert sql execute error: %s, %s", err.Error(), sql)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read last insert id from rs failed: %s", err.Error())
	}
	return id, nil
}
