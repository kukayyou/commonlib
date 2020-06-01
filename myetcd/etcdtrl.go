package myetcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/kukayyou/commonlib/mylog"
	"time"
)

func GetKey(etcdAddr string,key string)(value string){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		mylog.Error("connect etcd failed, err:%s", err.Error())
		return
	}

	mylog.Info("connect etcd success")
	defer cli.Close()

	//取值，设置超时为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		mylog.Error("get etcd key failed, key:%s, err:%s", key, err.Error())
		return
	}
	for _, ev := range resp.Kvs {
		value = string(ev.Value)
	}
	return
}
