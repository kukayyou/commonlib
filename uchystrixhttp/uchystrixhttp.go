package uchystrixhttp

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
	"go-agent/blueware/framework/beego"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/abursavich/nett"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/opentracing/opentracing-go"
)

const (
	Hystrix_Crit_Code int = 601
)

var (
	//HystrixEnable 启用hystrix
	HystrixEnable = false
	//HystrixStatEnable 启用hystrix 统计
	HystrixStatEnable = false
	//HystrixStatHTTPPort 启用hystrix 统计http端口
	HystrixStatHTTPPort = "8081"
	//Timeout 默认超时时间 ms
	Timeout = 5000
	//MaxConcurrent 最大并发数
	MaxConcurrent = 50
	//VolumeThreshold 最少统计请求数
	VolumeThreshold = 30
	//SleepWindow 心跳检测服务恢复时间 ms
	SleepWindow = 5000
	//ErrorPercentThreshold 错误百分比
	ErrorPercentThreshold = 50
	//HystrixConfig hystrix 命令配置
	HystrixConfig string
	//HystrixCommandName 默认CommandName
	HystrixCommandName = "default_command"
	//HystrixCommandConfigString CommandName规则
	HystrixCommandConfigString string
	//HystrixCommandConfig CommandName规则
	HystrixCommandConfig []map[string]string
	//UCCommandConfig CommandName规则, 当前配置
	UCCommandConfig map[string]UCConfig
	//ServerName 当前服务名称
	ServerName = ""
	//ServerIP 当前服务IP
	ServerIP = ""
	//DNSCacheResolver dns 缓存
	DNSCacheResolver = &nett.CacheResolver{TTL: 80 * time.Second}
)

// 当请求出现hystrix熔断失败时的回调处理函数
func HystrixFallback(e error) (int, []byte, *http.Response, error) {
	hystrixErr := fmt.Errorf("http hystrix error:%s", e.Error())
	uclog.Critical(hystrixErr.Error())
	return Hystrix_Crit_Code, nil, nil, hystrixErr
}

//UCConfig 当前包的配置
type UCConfig struct {
	hystrix.CommandConfig
	IngoreErrorCode []int `json:"ingore_error_code"` //要忽略的错误码
}

//UCRequest 请求
type UCRequest struct {
	commandName   string
	method        string
	url           string
	body          io.Reader
	header        map[string]string
	reqID         string
	cookies       []*http.Cookie
	tracerSpan    opentracing.Span
	tracerContext context.Context
	transport     *http.Transport
	dialer        *nett.Dialer
	basicAuthUser string
	basicAuthPass string
	statusCode    int
	resBody       []byte
	errMsg        error
}

//Init 初始化
func Init() error {
	if HystrixEnable == false {
		return nil
	}
	//默认配置
	hystrix.DefaultTimeout = Timeout
	hystrix.DefaultMaxConcurrent = MaxConcurrent
	hystrix.DefaultVolumeThreshold = VolumeThreshold
	hystrix.DefaultSleepWindow = SleepWindow
	hystrix.DefaultErrorPercentThreshold = ErrorPercentThreshold
	hystrix.ServerName = ServerName
	hystrix.ServerIP = ServerIP

	if err := json.Unmarshal([]byte(HystrixConfig), &UCCommandConfig); err != nil {
		return err
	}

	defaultCommandConfig := UCConfig{}
	defaultCommandConfig.Timeout = Timeout
	defaultCommandConfig.MaxConcurrentRequests = MaxConcurrent
	defaultCommandConfig.RequestVolumeThreshold = VolumeThreshold
	defaultCommandConfig.SleepWindow = SleepWindow
	defaultCommandConfig.ErrorPercentThreshold = ErrorPercentThreshold
	UCCommandConfig[HystrixCommandName] = defaultCommandConfig

	for commandName, config := range UCCommandConfig {
		hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
			Timeout:                config.Timeout,
			MaxConcurrentRequests:  config.MaxConcurrentRequests,
			RequestVolumeThreshold: config.RequestVolumeThreshold,
			SleepWindow:            config.SleepWindow,
			ErrorPercentThreshold:  config.ErrorPercentThreshold,
		})
	}
	if HystrixCommandConfigString != "" {
		if err := json.Unmarshal([]byte(HystrixCommandConfigString), &HystrixCommandConfig); err != nil {
			var oldConfig map[string]string
			if err := json.Unmarshal([]byte(HystrixCommandConfigString), &oldConfig); err != nil {
				return err
			}
			HystrixCommandConfig = append(HystrixCommandConfig, oldConfig)
		}
	}
	if HystrixStatEnable == true {
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrix.SetCallBackFunc(func(data hystrix.StreamCmdMetric) {
			metric, _ := json.Marshal(data)
			if data.CircuitBreakerOpen == true {
				uclog.Critical("hystrix circuit is open info :" + string(metric))
			}
			if data.ErrorPct >= 10 {
				uclog.Error("hystrix circuit ErrorPct 10%% info :" + string(metric))
			}
		}, func(data hystrix.StreamThreadPoolMetric) {
			metric, _ := json.Marshal(data)
			if (data.RollingCountThreadsExecuted*100)/data.CurrentMaximumPoolSize >= 95 {
				uclog.Critical("hystrix max concurrency 95% info :" + string(metric))
			}
		})
		hystrixStreamHandler.Start()
		go http.ListenAndServe(net.JoinHostPort("", HystrixStatHTTPPort), hystrixStreamHandler)
	}
	return nil
}

//New 创建请求
func New(method, url string, body io.Reader) *UCRequest {
	return &UCRequest{
		method:    method,
		url:       url,
		body:      body,
		transport: &http.Transport{},
		dialer: &nett.Dialer{
			Resolver: DNSCacheResolver,
			IPFilter: nett.DualStack,
		},
		header: make(map[string]string),
	}
}

//SetHeader 设置请求header
func (Req *UCRequest) SetHeader(header map[string]string) *UCRequest {
	for k, v := range header {
		Req.header[k] = v
	}
	for k, v := range header {
		if strings.ToLower(k) == "content-type" {
			Req.header["Content-Type"] = v
			break
		}
	}
	if _, ok := Req.header["Content-Type"]; !ok {
		Req.header["Content-Type"] = "application/json"
	}
	return Req
}

//SetCookies 设置cookie
func (Req *UCRequest) SetCookies(cookies []*http.Cookie) *UCRequest {
	if len(cookies) > 0 {
		Req.cookies = cookies
	}
	return Req
}

//SetCerts 设置证书
func (Req *UCRequest) SetCerts(certs []tls.Certificate) *UCRequest {
	Req.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	if len(certs) > 0 {
		Req.transport.TLSClientConfig.Certificates = certs
	}
	return Req
}

//SetProxy 设置代理
func (Req *UCRequest) SetProxy(proxy string) *UCRequest {
	Req.transport.Proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
	return Req
}

//SetAuth 设置验证
func (Req *UCRequest) SetAuth(basicAuthUser, basicAuthPass string) *UCRequest {
	Req.basicAuthUser = basicAuthUser
	Req.basicAuthPass = basicAuthPass
	return Req
}

//SetTimeout 这是超时
/*func (Req *UCRequest) SetTimeout(connTimeout, deadline time.Duration) *UCRequest {
	Req.dialer.Timeout = connTimeout
	Req.dialer.Deadline = time.Now().Add(deadline)
	return Req
}*/
func (Req *UCRequest) SetTimeout(connTimeout, deadline time.Duration) *UCRequest {
	//set dns dial
	Req.dialer.Timeout = connTimeout
	Req.dialer.Deadline = time.Now().Add(deadline)
	Req.transport.Dial = func(netw, addr string) (net.Conn, error) {
		c, err := Req.dialer.Dial(netw, addr)
		if err != nil {
			uclog.Warn("set dns error:%s,host:%s", err.Error(), addr)
			c, err := net.DialTimeout(netw, addr, connTimeout)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(deadline))
			return c, nil
		}
		c.SetDeadline(time.Now().Add(deadline))
		return c, nil

	}
	Req.transport.DisableKeepAlives = true
	return Req
}

//Request http请求
func (Req *UCRequest) Request(fallback func(error) (int, []byte, *http.Response, error)) (int, []byte, *http.Response, error) {
	Req.startSpan()
	defer Req.finish()
	var res *http.Response

	if HystrixEnable == true {
		commandName := ""
		for _, v := range HystrixCommandConfig {
			for tag, name := range v {
				if strings.Contains(Req.url, tag) {
					commandName = name
					break
				}
			}
			if commandName != "" {
				break
			}
		}
		txnData := oneapm_beego.GetTxnData()
		if commandName != "" {
			hystrix.Do(commandName, func() error {
				Req.statusCode, Req.resBody, res, Req.errMsg = Req.Request2_internal(txnData)
				//处理掉要忽略的错误码
				if len(UCCommandConfig[commandName].IngoreErrorCode) > 0 {
					for _, code := range UCCommandConfig[commandName].IngoreErrorCode {
						if code == Req.statusCode {
							uclog.Debug("IngoreErrorCode :%d", Req.statusCode)
							return nil
						}
					}
				}
				return Req.errMsg
			}, func(e error) error {
				Req.statusCode, Req.resBody, res, Req.errMsg = fallback(e)
				return nil
			})
			uclog.Debug("http request commandName:%s,url:%s", commandName, Req.url)
			return Req.statusCode, Req.resBody, res, Req.errMsg
		}
	}
	Req.statusCode, Req.resBody, res, Req.errMsg = Req.Request2()
	return Req.statusCode, Req.resBody, res, Req.errMsg
}

//Request2 http 请求
func (Req *UCRequest) Request2() (int, []byte, *http.Response, error) {
	return Req.Request2_internal(nil)
}

func (Req *UCRequest) Request2_internal(txnData *oneapm_beego.AgentBeegoData) (int, []byte, *http.Response, error) {
	starttime := time.Now()
	request, err := http.NewRequest(Req.method, Req.url, Req.body)

	if err != nil {
		httpErr := fmt.Errorf("construct http request failed, requrl = %s, err:%s", Req.url, err.Error())
		uclog.Critical(httpErr.Error())
		return -1, nil, nil, httpErr
	}
	if len(Req.header) > 0 {
		for k, v := range Req.header {
			request.Header.Add(k, v)
		}
		header, _ := json.Marshal(request.Header)
		uclog.Debug("%s request header value:%s", Req.reqID, string(header))
	}

	if Req.basicAuthUser != "" && Req.basicAuthPass != "" {
		uclog.Debug("%s request http auth user :%s, auth pass:%s", Req.reqID, Req.basicAuthUser, Req.basicAuthPass)
		request.SetBasicAuth(Req.basicAuthUser, Req.basicAuthPass)
	}
	for _, cookInfo := range Req.cookies {
		request.AddCookie(cookInfo)
	}
	//Req.transport.Dial = Req.dialer.Dial
	//Req.transport.DisableKeepAlives = true
	client := &http.Client{Transport: Req.transport}
	if txnData == nil {
		txnData = oneapm_beego.GetTxnData()
	}

	if txnData != nil {
		//client.Transport = blueware.NewRoundTripper(txnData.GetTxn(), Req.transport)
		uclog.Debug("Add header Nodeguid")

		//增加header选项
		request.Header.Add("Nodeguid", txnData.GetGuid(Req.url))
		request.Header.Add("Nodetrip", txnData.GetTripid(Req.url))
	}

	response, err := client.Do(request)
	if err != nil {
		httpErr := fmt.Errorf("http request failed, url: %s, error:%s", Req.url, err.Error())
		uclog.Critical(httpErr.Error())
		return -1, nil, nil, httpErr
	}

	if response != nil {
		defer response.Body.Close()
		endtime := time.Since(starttime)
		uclog.Info("%s http request method:%s, request url:%s,request duration: %s", Req.reqID, Req.method, Req.url, endtime.String())
		var respBody []byte
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(response.Body)
			defer reader.Close()
			respBody, _ = ioutil.ReadAll(reader)
		default:
			respBody, _ = ioutil.ReadAll(response.Body)
		}
		if response.StatusCode > 300 {
			httpErr := fmt.Errorf("http request failed, http method:%s,status:%d,url:%s,response:%s", Req.method, response.StatusCode, Req.url, string(respBody))
			// UMS 404属于数据查询不到时的正常业务逻辑
			if response.StatusCode != 404 {
				uclog.Critical(httpErr.Error())
			}

			return response.StatusCode, respBody, response, httpErr
		}
		return response.StatusCode, respBody, response, nil
	}
	return -1, nil, nil, fmt.Errorf("%s to %s, empty response", Req.method, Req.url)
}
