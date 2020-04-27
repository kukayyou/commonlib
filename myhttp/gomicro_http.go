package myhttp

import (
	"context"
	"fmt"
	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	microhttp "github.com/micro/go-plugins/client/http"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
)

func RequestWithHytrix(serverName, url string, req interface{})map[string]interface{}{
	consulReg := consul.NewRegistry(
		registry.Addrs("192.168.109.131:8500"),
	)

	microselector := selector.NewSelector(
		selector.Registry(consulReg),              //传入consul注册
		selector.SetStrategy(selector.RoundRobin), //指定查询机制
	)
	microClient := microhttp.NewClient(
		client.Selector(microselector),
		client.ContentType("application/json"),
		client.Wrap(hystrix.NewClientWrapper()), //熔断操作
	)
	hystrixGo.DefaultSleepWindow = 5000//重试时间窗口
	hystrixGo.DefaultTimeout = 5000//默认超时时间
	hystrixGo.DefaultVolumeThreshold = 2//默认最大失败次数

	reqInfo := microClient.NewRequest(serverName, url, req)
	var resp map[string]interface{}

	err := microClient.Call(context.Background(), reqInfo, &resp)
	if err == nil {
		fmt.Println(resp)
	}

	return  resp
}