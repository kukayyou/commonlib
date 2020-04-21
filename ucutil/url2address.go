package ucutil

import (
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

var url2AddressMap map[string]map[int64]string

func InitUrl2Address(datasrc string) {

	url2AddressMap = make(map[string]map[int64]string, 0)

	if len(datasrc) > 0 {
		o := orm.NewOrm()
		o.Using("iptables")

		sql := "select url, type, ip_address from uc_url2ipaddress"
		uclog.Debug("get url2ipaddress from db sql :%s", sql)

		var maps []orm.Params
		num, err := o.Raw(sql).Values(&maps)

		if err != nil {
			uclog.Error("search url2ipaddress from db failed, sql: %s, err: %s", sql, err.Error())
			return
		}

		if num < 1 {
			uclog.Info("search url2ipaddress result is empty")
			return
		}

		uclog.Debug("======================== URL2ADDRESS ========================")
		for _, row := range maps {
			url := ToString(row["url"])
			ipType := ToInt64(row["type"], 0)
			ipAddress := ToString(row["ip_address"])

			if len(url2AddressMap[url]) == 0 {
				url2AddressMap[url] = make(map[int64]string, 0)
			}

			url2AddressMap[url][ipType] = ipAddress

		}

		uclog.Debug("====================== END URL2ADDRESS ======================")

	}
}

func GetAddrByURL(url string, ipType int64) string {

	if len(url2AddressMap) == 0 {
		uclog.Warn("url2AddressMap has no data")
		return ""
	}

	if v, ok := url2AddressMap[url]; ok {
		if address, ok1 := v[ipType]; ok1 {
			uclog.Info("find the address by url:%s and type:%d and address:%s", url, ipType, address)
			return address
		}
	}

	return ""
}

func ParseUrl(strUrl string) (path string, host string, port int, err error) {
	if strUrl == "" {
		err = fmt.Errorf("Invalid params of strUrl empty")
		return
	}
	u, err := url.Parse(strUrl)
	if err != nil {
		return
	}

	path = u.Path
	h := strings.Split(u.Host, ":")
	host = h[0]
	if len(h) >= 2 {
		port, _ = strconv.Atoi(h[1])
	}
	return
}
