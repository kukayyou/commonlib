package ucredis

import (
	"github.com/garyburd/redigo/redis"
)

// redis单条命令的执行状态
const (
	COMMANDSTATUS_NOTEXEC = 0 // 未执行
	COMMANDSTATUS_EXECED  = 1 // 已执行
)

type slotRange struct {
	min uint16
	max uint16
}

// 表示redis单个命令
type UcRedisCommand struct {
	UserData interface{} // user data
	Key      string
	Cmd      string
	Params   []interface{}
	Reply    interface{}
	Status   int
}

// 表示cluster中的某个节点
type ucRedisNode struct {
	conn   *redis.Conn       //redis链接
	slots  []slotRange       //hash slot范围
	cmdLst []*UcRedisCommand // support for pipeline
}

func (n *ucRedisNode) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return (*n.conn).Do(commandName, args...)
}
func (n *ucRedisNode) HasPipeCommand() bool {
	return len(n.cmdLst) > 0
}
func (n *ucRedisNode) ClearPipeCommand() {
	n.cmdLst = make([]*UcRedisCommand, 0)
}
