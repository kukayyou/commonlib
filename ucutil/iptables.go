package ucutil

import (
	"encoding/binary"
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
	"net"

	"github.com/astaxie/beego/orm"
)

const (
	Domestic = iota
	Oversea
	Mobile
	HongKang
)

var sIpTableMap map[uint32]uint32 // source ip ---> end ip
var sIpTypeMap map[uint32]int64   // source ip ---> type
var sIpSlice []uint32             // source ip slice

func Init(datasrc string) {

	sIpTableMap = make(map[uint32]uint32, 0)
	sIpTypeMap = make(map[uint32]int64, 0)
	sIpSlice = make([]uint32, 0)

	if len(datasrc) > 0 {
		o := orm.NewOrm()
		o.Using("iptables")

		sql := "select StartIp, EndIp, Type from uc_ipaddress"
		uclog.Debug("get iptables from db sql :%s", sql)

		var maps []orm.Params
		num, err := o.Raw(sql).Values(&maps)

		if err != nil {
			uclog.Error("search iptables from db failed, sql: %s, err: %s", sql, err.Error())
			return
		}

		if num < 1 {
			uclog.Info("search iptables result is empty")
			return
		}

		uclog.Debug("======================== IPTABLES ========================")
		for _, row := range maps {
			startIpStr := ToString(row["StartIp"])
			endIpStr := ToString(row["EndIp"])
			ipType := ToInt64(row["Type"], 0)

			startIp, err := IpaddrToInt(startIpStr)
			if err != nil {
				uclog.Error("ipaddrToInt failed. source ip address: %s, err: %s", startIpStr, err.Error())
				continue
			}

			endIp, err := IpaddrToInt(endIpStr)
			if err != nil {
				uclog.Error("ipaddrToInt failed. end ip address: %s, err: %s", endIpStr, err.Error())
				continue
			}

			sIpTableMap[startIp] = endIp
			sIpTypeMap[startIp] = ipType
			sIpSlice = append(sIpSlice, startIp)

		}
		uclog.Debug("====================== END IPTABLES ======================")

		SortUint32(sIpSlice)
	}
}

func IpaddrToInt(ip string) (uint32, error) {

	if len(ip) == 0 {
		return 0, fmt.Errorf("ip address is empty")
	}

	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return 0, fmt.Errorf("%v is not an IPv4 address", trial)
	}
	return binary.BigEndian.Uint32(trial.To4()), nil
}

func BinarySearch(s []uint32, k uint32) (int, int, int) {

	if s[0] > k {
		return -1, 0, 0
	}

	if s[len(s)-1] < k {
		return -1, len(s) - 1, len(s) - 1
	}

	lo, p_lo, hi, p_hi := 0, 0, len(s)-1, len(s)-1

	for lo <= hi {

		p_lo = lo
		p_hi = hi

		m := (lo + hi) >> 1
		if s[m] < k {
			lo = m + 1
		} else if s[m] > k {
			hi = m - 1
		} else {
			return m, p_lo, p_hi
		}
	}

	return -1, p_lo, p_hi
}

//true means ip source from overseas, false means ip source from domestic city
func CheckIp(ip string) int64 {

	if len(ip) == 0 || len(sIpTableMap) == 0 || len(sIpSlice) == 0 {
		uclog.Info("ip, sIpTableMap, sIpSlice is empty")
		return Domestic
	}

	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		uclog.Error("%v is not an IPv4 address", trial)
		return Domestic
	}

	u_ip, err := IpaddrToInt(ip)
	if err != nil {
		uclog.Error("parse ip address: %s failed, error: %s.", ip, err.Error())
		return Domestic
	}

	if _, ok := sIpTableMap[u_ip]; ok {
		uclog.Info("ip address hit in sIpTableMap directly. ip address: %s and ip type: %d", ip, sIpTypeMap[u_ip])
		return sIpTypeMap[u_ip]
	}

	// search in source ip slice
	_, l, _ := BinarySearch(sIpSlice, u_ip)
	var sIp uint32
	if l > 0 {
		sIp = sIpSlice[l-1]
		if eIp, ok := sIpTableMap[sIp]; ok {
			if u_ip >= sIp && u_ip <= eIp {
				uclog.Info("hit ip address:%s in [%d, %d] range and ip type: %d in range", ip, sIp, eIp, sIpTypeMap[sIp])
				return sIpTypeMap[sIp]
			}
		}
	}
	sIp = sIpSlice[l]
	if eIp, ok := sIpTableMap[sIp]; ok {
		if u_ip >= sIp && u_ip <= eIp {
			uclog.Info("hit ip address:%s in [%d, %d] range and ip type: %d in edge", ip, sIp, eIp, sIpTypeMap[sIp])
			return sIpTypeMap[sIp]
		}
	}

	uclog.Info("we donot find this ip address:%s in iptables", ip)
	return Domestic

}
