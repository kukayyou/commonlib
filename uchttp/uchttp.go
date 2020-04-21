package uchttp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gnetis.com/golang/core/golib/uchystrixhttp"
	"gnetis.com/golang/core/golib/uclog"
)

var (
	enableEncode  bool
	basicAuthUser = ""
	basicAuthPass = ""
)

const (
	HTTP_TIMEOUT_CONNECT  time.Duration = 3 * time.Second
	HTTP_TIMEOUT_DEADLINE time.Duration = 5 * time.Second
)

func ReplaceHost(rawUrl, host string, prefixs ...string) string {
	if strings.Contains(host, "://") {
		hostUrl, _ := url.Parse(host)
		if hostUrl != nil {
			host = hostUrl.Host
		}
	}
	if host == "" {
		return rawUrl
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	u.Host = host

	if len(prefixs) > 0 && len(prefixs[0]) > 0 {
		prefix := prefixs[0]
		prefix = strings.TrimSuffix(prefix, "/")
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		u.Path = prefix + u.Path
	}
	return u.String()
}

func MakeHttpLink(base, path string) string {
	link := MakeLink(base, path)
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "http://" + link
	}
	return link
}

func MakeLink(base, path string) string {
	if strings.HasSuffix(base, "/") {
		base = strings.TrimRight(base, "/")
	}
	if strings.HasPrefix(path, "/") {
		path = strings.TrimLeft(path, "/")
	}
	return base + "/" + path
}

func InitHttpHeader() {
	enableEncode = true
}

func InitHttpAuthData(authUser, authPass string) {
	basicAuthUser = authUser
	basicAuthPass = authPass
}

func CookieRequestWithRetry(method string, requrl string, body string) (int, []byte, error) {
	return RequestCookie(method, requrl, body, HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE, true, nil)
}
func RequestCookie(method string, requrl string, body string, connTimeout, deadline time.Duration, retry bool, cookies []*http.Cookie) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, strings.NewReader(body))
	request.SetTimeout(connTimeout, deadline)
	request.SetCookies(cookies)
	request.SetCerts(make([]tls.Certificate, 0, 0))
	header := map[string]string{"Content-Type": "application/json"}
	if enableEncode {
		header["Accept-Encoding"] = "gzip"
	}
	request.SetHeader(header)
	var (
		statusCode int
		resBody    []byte
		res        *http.Response
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, res, err = request.Request(fallback)
	if statusCode == 444 && retry {
		j, _ := json.Marshal(res.Cookies())
		uclog.Info("http response 444, retry with cookies:%s, body:%+v", string(j), body)
		// request.SetCookies(res.Cookies())
		// statusCode, resBody, _, err = request.Request(fallback)
		return RequestCookie(method, requrl, body, connTimeout, deadline, false, res.Cookies())
	}
	return statusCode, resBody, err
}

func Request(method string, requrl string, body io.Reader) (int, []byte, error) {
	return RequestTimeout(method, requrl, body, HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
}

func RequestWithHeader(method string, requrl string, body io.Reader, header map[string]string) (int, []byte, error) {
	return RequestTimeoutWithHeader(method, requrl, body, HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE, header)
}

func RequestTimeout(method string, requrl string, body io.Reader, connTimeout, deadline time.Duration) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(connTimeout, deadline)
	header := map[string]string{"Content-Type": "application/json"}
	if enableEncode {
		header["Accept-Encoding"] = "gzip"
	}
	request.SetHeader(header)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func SSLRequest(method string, requrl string, body io.Reader, cert ...tls.Certificate) (int, []byte, error) {
	return SSLRequestTimeout(method, requrl, body, HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE, cert...)
}
func SSLRequestTimeout(method string, requrl string, body io.Reader, connTimeout, deadline time.Duration, certs ...tls.Certificate) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(connTimeout, deadline)
	header := map[string]string{"Content-Type": "application/json"}
	if enableEncode {
		header["Accept-Encoding"] = "gzip"
	}
	request.SetHeader(header)
	request.SetCerts(certs)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func RequestForm(method string, requrl string, body io.Reader) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
	request.SetHeader(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	request.SetAuth(basicAuthUser, basicAuthPass)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func RequestXML(method string, requrl string, body io.Reader) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
	request.SetHeader(map[string]string{"Content-Type": "application/xml"})
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func RequestXMLWithHeader(method string, requrl string, body io.Reader, httpHeader map[string]string) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
	if len(httpHeader) > 0 {
		httpHeader["Content-Type"] = "application/xml"
	}
	request.SetHeader(httpHeader)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func RequestFormWithProxy(method string, requrl string, body io.Reader, proxy string) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
	request.SetHeader(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	request.SetProxy(proxy)
	request.SetAuth(basicAuthUser, basicAuthPass)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func ProxyRequestWithHeader(proxyUrl, method string, requrl string, body io.Reader, httpHeader map[string]string) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
	if len(httpHeader) > 0 {
		httpHeader["Content-Type"] = "application/json"
	}
	request.SetHeader(httpHeader)
	request.SetProxy(proxyUrl)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func ProxyRequest(proxyUrl, method string, requrl string, body io.Reader) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(HTTP_TIMEOUT_CONNECT, HTTP_TIMEOUT_DEADLINE)
	request.SetHeader(map[string]string{"Content-Type": "application/json"})
	request.SetProxy(proxyUrl)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

func RequestTimeoutWithHeader(method string, requrl string, body io.Reader, connTimeout, deadline time.Duration, header map[string]string) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(connTimeout, deadline)
	for k, v := range header {
		if strings.ToLower(k) == "content-type" {
			delete(header, k)
			header["Content-Type"] = v
		}
	}
	if len(header) > 0 {
		if _, ok := header["Content-Type"]; !ok {
			header["Content-Type"] = "application/json"
		}
	}

	request.SetHeader(header)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)

	return statusCode, resBody, err
}

func ProxyRequestTimeoutWithHeader(proxyUrl, method string, requrl string, body io.Reader, connTimeout, deadline time.Duration, httpHeader map[string]string) (int, []byte, error) {
	request := uchystrixhttp.New(method, requrl, body)
	request.SetTimeout(connTimeout, deadline)
	if len(httpHeader) > 0 {
		httpHeader["Content-Type"] = "application/json"
	}
	request.SetHeader(httpHeader)
	request.SetProxy(proxyUrl)
	var (
		statusCode int
		resBody    []byte
		err        error
		fallback   = uchystrixhttp.HystrixFallback
	)
	statusCode, resBody, _, err = request.Request(fallback)
	return statusCode, resBody, err
}

type RequestOption struct {
	logger        *uclog.UcLog      // 日志记录器
	logBodyLength int               // 输出消息体最大长度，默认为512
	dialTimeout   time.Duration     // 连接超时, 默认为3s
	callTimeout   time.Duration     // 读写超时, 默认为5s
	header        map[string]string // 请求头， 默认为"Content-Type": "application/json"
	proxyUrl      string            // 代理服务器地址
}

func NewRequestOption() *RequestOption {
	opt := &RequestOption{}
	opt.logger = nil
	opt.logBodyLength = 1024
	opt.proxyUrl = ""
	opt.dialTimeout = HTTP_TIMEOUT_CONNECT
	opt.callTimeout = HTTP_TIMEOUT_DEADLINE
	opt.header = map[string]string{"Content-Type": "application/json"}
	return opt
}

func GetOrNewRequestOption(opts ...*RequestOption) *RequestOption {
	var opt *RequestOption
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = NewRequestOption()
	}
	return opt
}

func (r *RequestOption) Logger(logger *uclog.UcLog) *RequestOption {
	if logger != nil {
		r.logger = logger
	}
	return r
}

func (r *RequestOption) LogBodyLength(logBodyLength int) *RequestOption {
	if logBodyLength > 0 {
		r.logBodyLength = logBodyLength
	}
	return r
}

// 3s = 3 * 1000 * 1000 * 1000
func (r *RequestOption) Timeout(dialTimeout, callTimeout int) *RequestOption {
	if dialTimeout > 0 {
		r.dialTimeout = time.Duration(dialTimeout)
	}
	if callTimeout > 0 {
		r.callTimeout = time.Duration(callTimeout)
	}
	return r
}

func (r *RequestOption) Header(header map[string]string) *RequestOption {
	if header != nil {
		r.header = header
	}
	return r
}

func (r *RequestOption) Proxy(proxyUrl string) *RequestOption {
	if len(proxyUrl) > 0 {
		r.proxyUrl = proxyUrl
	}
	return r
}

/*
1, 统一输出请求和响应参数用于定位问题
2, 输出请求耗时用于分析性能
3, 作为一行输出减少日志量
4, 网络失败时输出Critical日志用于报警
5, 去除响应包含的换行方便搜索
6, 如果请求体太大，对内容进行截断
*/
func RequestWithMonitor(method string, reqUrl string, body string, opts ...*RequestOption) (int, []byte, error) {
	var (
		alarm     = false
		starttime = time.Now()
		logBody   = body
		code      int
		response  []byte
		err       error
	)
	opt := GetOrNewRequestOption(opts...)
	if end := opt.logBodyLength; len(logBody) > end {
		logBody = logBody[0:end] + "..."
	}
	requestInfo := fmt.Sprintf("method:%s, reqUrl:%s, header:%v, body:%s,", method, reqUrl, opt.header, logBody)
	defer func() {
		duration := time.Since(starttime).String()
		content := fmt.Sprintf("requestInfo duration:%s, %s", duration, requestInfo)
		if alarm {
			if opt.logger != nil {
				opt.logger.Log_Error("%s", content)
				opt.logger.Log_Critical("%s", content)
			} else {
				uclog.Error("%s", content)
				uclog.Critical("%s", content)
			}
		} else {
			if opt.logger != nil {
				opt.logger.Log_Info("%s", content)
			} else {
				uclog.Info("%s", content)
			}
		}
	}()

	if opt.proxyUrl != "" {
		code, response, err = ProxyRequestTimeoutWithHeader(opt.proxyUrl, method, reqUrl, strings.NewReader(body), opt.dialTimeout, opt.callTimeout, opt.header)
	} else {
		code, response, err = RequestTimeoutWithHeader(method, reqUrl, strings.NewReader(body), opt.dialTimeout, opt.callTimeout, opt.header)
	}

	if err != nil {
		requestInfo += fmt.Sprintf("error:%s", err.Error())
		alarm = true
	} else {
		resp := strings.Replace(string(response), "\n", "", -1)
		requestInfo += fmt.Sprintf("code:%d, resp:%s", code, resp)
		if code > 300 {
			alarm = true
		}
	}
	return code, response, err
}
