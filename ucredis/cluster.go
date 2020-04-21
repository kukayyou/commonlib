package ucredis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gnetis.com/golang/core/golib/uclog"
	"strconv"
	"strings"
	"time"
)

var (
	crc16tab = [256]uint16{
		0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
		0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
		0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
		0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
		0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
		0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
		0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
		0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
		0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
		0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
		0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
		0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
		0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
		0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
		0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
		0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
		0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
		0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
		0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
		0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
		0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
		0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
		0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
		0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
		0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
		0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
		0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
		0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
		0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
		0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
		0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
		0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0}
)

func NewUcRedis() *UcRedisCluster {
	r := UcRedisCluster{}
	r.addrs = make([]string, 0)
	r.nodes = make([]ucRedisNode, 0)
	return &r
}

func GetClusterAddress(addrs []string) ([]string, error) {
	var err error
	var conn redis.Conn
	for i := 0; i < len(addrs); i++ {
		conn, err = redis.DialTimeout("tcp", addrs[i], 5*time.Second, 5*time.Second, 5*time.Second)
		if err != nil {
			uclog.Error("Connect to redis fail, %s", err.Error())
			continue
		}
		break
	}

	reply, err := conn.Do("cluster", "nodes")
	ret, err := redis.String(reply, err)
	if err != nil {
		uclog.Error("exec [cluster nodes] fail, %s", err.Error())
		return nil, err
	}

	r := make([]string, 0)

	lines := strings.Split(ret, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		if !strings.Contains(fields[2], "master") ||
			strings.Contains(fields[2], "fail") {
			continue
		}

		r = append(r, fields[1])
	}

	return r, nil
}

// 表示cluster
type UcRedisCluster struct {
	addrs    []string
	nodes    []ucRedisNode
	password string
	pipeCmd  []UcRedisCommand
}

func (cluster *UcRedisCluster) Nodes() []ucRedisNode {
	return cluster.nodes
}

func (cluster *UcRedisCluster) AddRedisServer(addr string) {
	cluster.addrs = append(cluster.addrs, addr)
}

func (cluster *UcRedisCluster) SetPassword(password string) {
	cluster.password = password
	uclog.Info("set password: %s", password)
}

func (cluster *UcRedisCluster) Dial() bool {
	return cluster.reconnect()
}

func (cluster *UcRedisCluster) Close() {
	for _, conns := range cluster.nodes {
		(*conns.conn).Close()
	}
	cluster.nodes = make([]ucRedisNode, 0)
}

func (cluster *UcRedisCluster) Do(key string, cmd string, args ...interface{}) (r interface{}, err error) {
	r, err = cluster.excute(true, key, cmd, args...)
	if err != nil {
		uclog.Critical("redis error:%s", err.Error())
		return
	}

	return
}

func (cluster *UcRedisCluster) Send(key string, cmd string, args ...interface{}) error {
	return cluster.send(true, key, cmd, args...)
}

func (cluster *UcRedisCluster) Eval(key string, script *redis.Script, args ...interface{}) (interface{}, error) {
	return cluster.eval(true, key, script, args...)
}

func (cluster *UcRedisCluster) reconnect() bool {
	cluster.Close()
	var c *redis.Conn
	for i := 0; i < len(cluster.addrs); i++ {
		conn, err := redis.DialTimeout("tcp", cluster.addrs[i], 5*time.Second, 5*time.Second, 5*time.Second)
		if err != nil {
			uclog.Error("Connect to redis fail, %s", err.Error())
			continue
		}

		// 如果需要密码验证，校验密码
		if len(cluster.password) > 0 {
			if _, err := conn.Do("AUTH", cluster.password); err != nil {
				uclog.Error("auth fail by password: %s, addr: %s, err: %s", cluster.password, cluster.addrs[i], err.Error())
				conn.Close()
				continue
			}
		}

		c = &conn
		break
	}

	if c == nil {
		return false
	}
	defer (*c).Close()
	return cluster.initCluster(c)
}

func (cluster *UcRedisCluster) initCluster(c *redis.Conn) bool {
	cluster.Close()

	reply, err := (*c).Do("cluster", "nodes")
	ret, err := redis.String(reply, err)
	if err != nil {
		uclog.Error("initCluster fail, %s", err.Error())
		return false
	}

	lines := strings.Split(ret, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		if !strings.Contains(fields[2], "master") ||
			strings.Contains(fields[2], "fail") {
			continue
		}

		addr := fields[1]
		uclog.Info("connect to redis host: %s, password: %s", addr, cluster.password)
		conn, err := redis.DialTimeout("tcp", addr, 5*time.Second, 5*time.Second, 5*time.Second)
		if err != nil {
			uclog.Error("Connect to redis fail, %s", err.Error())
			cluster.Close()
			return false
		}

		// 如果需要密码验证，校验密码
		if len(cluster.password) > 0 {
			if _, err := conn.Do("AUTH", cluster.password); err != nil {
				uclog.Error("auth fail by password: %s, addr: %s, err: %s", cluster.password, addr, err.Error())
				cluster.Close()
				return false
			}
		}

		node := ucRedisNode{}
		node.slots = make([]slotRange, 0)
		for j := 8; j < len(fields); j++ {
			subFields := strings.Split(fields[j], "-")
			var slot slotRange
			var min, max int
			if len(subFields) == 2 {
				min, _ = strconv.Atoi(subFields[0])
				max, _ = strconv.Atoi(subFields[1])
			} else {
				min, _ = strconv.Atoi(subFields[0])
				max, _ = strconv.Atoi(subFields[0])
			}
			slot.min = uint16(min)
			slot.max = uint16(max)
			node.slots = append(node.slots, slot)
			node.conn = &conn
		}
		cluster.nodes = append(cluster.nodes, node)
	}

	return true
}

func (cluster *UcRedisCluster) excute(retry bool, key string, cmd string, args ...interface{}) (interface{}, error) {
	node := cluster.getNode(key)
	if node == nil || (*node).conn == nil {
		cluster.reconnect()
		return nil, fmt.Errorf("redis cluster error!")
	}

	reply, err := (*(*node).conn).Do(cmd, args...)
	if err != nil {
		if retry {
			uclog.Warn("ucReis do fail,  key:%s, cmd:%s, error:%s", key, cmd, err.Error())
			cluster.reconnect()
			return cluster.excute(false, key, cmd, args...)
		} else {
			uclog.Error("ucReis retry do fail,  key:%s, cmd:%s, error:%s", key, cmd, err.Error())
		}
	}
	return reply, err
}

func (cluster *UcRedisCluster) send(retry bool, key string, cmd string, args ...interface{}) error {
	node := cluster.getNode(key)
	if node == nil || (*node).conn == nil {
		cluster.reconnect()
		return fmt.Errorf("redis cluster error!")
	}

	err := (*(*node).conn).Send(cmd, args...)
	if err != nil {
		if retry {
			uclog.Warn("ucReis do fail,  key:%s, cmd:%s, error:%s", key, cmd, err.Error())
			cluster.reconnect()
			return cluster.send(false, key, cmd, args...)
		} else {
			uclog.Error("ucReis retry do fail,  key:%s, cmd:%s, error:%s", key, cmd, err.Error())
		}
	}
	return err
}

func (cluster *UcRedisCluster) eval(retry bool, key string, script *redis.Script, args ...interface{}) (interface{}, error) {
	node := cluster.getNode(key)
	if node == nil || (*node).conn == nil {
		cluster.reconnect()
		return nil, fmt.Errorf("redis cluster error!")
	}

	reply, err := script.Do(*(*node).conn, args...)
	if err != nil {
		if retry {
			uclog.Warn("ucReis do fail,  key:%s, cmd:%s, error:%s", key, "eval", err.Error())
			cluster.reconnect()
			return cluster.eval(false, key, script, args...)
		} else {
			uclog.Error("ucReis retry do fail,  key:%s, cmd:%s, error:%s", key, "eval", err.Error())
		}
	}
	return reply, err
}

func (cluster *UcRedisCluster) hashSlot(key string) uint16 {
	crc := uint16(0)
	for i := 0; i < len(key); i++ {
		crc = (crc << 8) ^ crc16tab[((crc>>8)^uint16(key[i]))&0x00FF]
	}

	return crc % 16384
}

func (cluster *UcRedisCluster) getNode(key string) *ucRedisNode {
	slot := cluster.hashSlot(key)
	for i := 0; i < len(cluster.nodes); i++ {
		for j := 0; j < len(cluster.nodes[i].slots); j++ {
			if slot >= cluster.nodes[i].slots[j].min &&
				slot <= cluster.nodes[i].slots[j].max {
				return &(cluster.nodes[i])
			}
		}
	}

	return nil
}

// support for pipeline
func (cluster *UcRedisCluster) Pipe(ud interface{}, key string, cmd string, args ...interface{}) {
	rcmd := UcRedisCommand{UserData: ud, Key: key, Cmd: cmd, Params: args}
	cluster.pipeCmd = append(cluster.pipeCmd, rcmd)
}

func (cluster *UcRedisCluster) Commit() (reply []UcRedisCommand, err error) {

	defer func() {

		reply = cluster.pipeCmd
		cluster.pipeCmd = make([]UcRedisCommand, 0)

		if r := recover(); r != nil {
			s := fmt.Sprintln("redis pipeline panic:", r)
			uclog.Error(s)
			err = fmt.Errorf("%s", s)
		}
	}()

	err = cluster.commitAll()
	if err != nil {
		cluster.reconnect()
		err = cluster.commitAll()
	}
	if err != nil {
		uclog.Error("redis pipeline commit error:%s", err.Error())
	}

	return
}

func (cluster *UcRedisCluster) commitAll() error {

	var saved_err error = nil

	saved_err = cluster.dispatchCommand()
	if saved_err != nil {
		return saved_err
	}

	// 循环请求cluster中的各个redis实例
	for i := 0; i < len(cluster.nodes); i++ {
		if cluster.nodes[i].HasPipeCommand() {
			err := cluster.commitNode(i)
			if err != nil {
				saved_err = err
				uclog.Warn("redis pipeline commit node failed, %s", err.Error())
			}
			cluster.nodes[i].ClearPipeCommand()
		}
	}

	return saved_err
}

func (cluster *UcRedisCluster) commitNode(index int) (saved_err error) {

	node := &(cluster.nodes[index])
	cmdList := node.cmdLst
	conn := node.conn

	for _, redisCmd := range cmdList {
		saved_err = (*conn).Send(redisCmd.Cmd, redisCmd.Params...)
		if saved_err != nil {
			return
		}
	}

	// 发送缓冲区的命令到redis，接收结果并解析
	saved_err = (*conn).Flush()
	if saved_err != nil {
		uclog.Warn("commitNode flush to server error:%s", saved_err.Error())
		return
	}

	for i := 0; i < len(cmdList); i++ {
		reply, err := (*conn).Receive()
		if err != nil {
			if saved_err == nil {
				saved_err = err
			}
			uclog.Warn("commitNode receive from server error:%s", saved_err.Error())
			continue
		}
		node.cmdLst[i].Reply = reply
		node.cmdLst[i].Status = COMMANDSTATUS_EXECED
	}

	return
}

func (cluster *UcRedisCluster) dispatchCommand() error {
	for i := 0; i < len(cluster.pipeCmd); i++ {
		cmd := &cluster.pipeCmd[i]

		if cmd.Status == COMMANDSTATUS_EXECED {
			continue
		}
		redisNode := cluster.getNode(cmd.Key)
		if redisNode == nil {
			return fmt.Errorf("cannot dispatch command to redis node")
		}

		(*redisNode).cmdLst = append((*redisNode).cmdLst, cmd)
	}

	return nil
}
