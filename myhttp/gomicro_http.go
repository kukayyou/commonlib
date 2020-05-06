package myhttp

import (
	"context"
	"fmt"
	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	microhttp "github.com/micro/go-plugins/client/http"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
)

var (
	ConsulAddr string//consul地址：ip+port
	EtcdAddr string//consul地址：ip+port
	DefaultSleepWindow int = 5000//重试时间窗口
	DefaultTimeOut int = 5000//默认超时时间
	DefaultVolumeThreshold int = 2//默认最大失败次数
	RegistryType int = 0//0:etcd ,1:consul
)

func RequestWithHytrix(serverName, url string, req interface{})map[string]interface{}{
	var reg registry.Registry
	switch RegistryType {
	case 0:
		reg = etcd.NewRegistry(
			registry.Addrs(EtcdAddr),
		)
	case 1:
		reg = consul.NewRegistry(
			registry.Addrs(ConsulAddr),
		)
	default:
	}

	microselector := selector.NewSelector(
		selector.Registry(reg),              //传入consul注册
		selector.SetStrategy(selector.RoundRobin), //指定查询机制
	)
	microClient := microhttp.NewClient(
		client.Selector(microselector),
		client.ContentType("application/json"),
		client.Wrap(hystrix.NewClientWrapper()), //熔断操作
	)
	hystrixGo.DefaultSleepWindow = DefaultSleepWindow//重试时间窗口
	hystrixGo.DefaultTimeout = DefaultTimeOut//默认超时时间
	hystrixGo.DefaultVolumeThreshold = DefaultVolumeThreshold//默认最大失败次数

	reqInfo := microClient.NewRequest(serverName, url, req)
	var resp map[string]interface{}

	err := microClient.Call(context.Background(), reqInfo, &resp)
	if err == nil {
		fmt.Println(resp)
	}

	return  resp
}