package myetcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/kukayyou/commonlib/mylog"
)

func GetKey(etcdAddr string,key string)(value string){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
	})
	if err != nil {
		mylog.Error("connect etcd failed, err:%s", err.Error())
		return
	}

	mylog.Info("connect etcd success")
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	resp, err := kv.Get(context.TODO(), key)
	if err != nil {
		mylog.Error("get etcd key failed, key:%s, err:%s", key, err.Error())
		return
	}
	for _, ev := range resp.Kvs {
		value = string(ev.Value)
	}
	return
}