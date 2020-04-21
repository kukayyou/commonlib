package uchttp

import (
	"crypto/tls"
	"gnetis.com/golang/core/golib/uchystrixhttp"
	"gnetis.com/golang/core/golib/uclog"
	"net/http"
	"strings"
	"testing"
	"time"
)

func InitTest() {
	uclog.Initialize("log", `D:\Go\WorkSpace\src\golib\uchystrixhttp\a.log`, "debug")
}

func TestMakeLink(t *testing.T) {
	if MakeLink("http://www.quanshi.com", "aaa/aaaa/a") != "http://www.quanshi.com/aaa/aaaa/a" {
		t.Error("TestMakeLink error")
	}
}

func TestRequest1(t *testing.T) {
	InitTest()
	_, b, err := Request("GET", "http://www.baidu.com", nil)
	t.Log(string(b))
	t.Error(err)
}

func TestRequest2(t *testing.T) {
	InitTest()
	_, b, err := CookieRequestWithRetry("POST", "http://www.baidu.com", "")
	t.Log(string(b))
	t.Error(err)
}

func TestRequest3(t *testing.T) {
	InitTest()
	_, b, err := RequestCookie("GET", "http://family.quanshi.com/new_index.php?m=&&m=Hr&c=index&a=index", "", time.Duration(1000)*time.Millisecond, time.Duration(1000)*time.Millisecond, true, []*http.Cookie{&http.Cookie{Name: "PHPSESSID", Value: "79d2rliq62lkv3cqeerda4pok2"}})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest4(t *testing.T) {
	InitTest()
	_, b, err := RequestWithHeader("GET", "http://www.baidu.com", nil, map[string]string{"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest5(t *testing.T) {
	InitTest()
	_, b, err := RequestTimeout("GET", "http://family.quanshi.com/new_index.php?m=&&m=Hr&c=index&a=index", nil, time.Duration(1000)*time.Millisecond, time.Duration(1000)*time.Millisecond)
	t.Log(string(b))
	t.Error(err)
}

func TestRequest6(t *testing.T) {
	InitTest()
	cliCrt, err := tls.LoadX509KeyPair(`D:\Go\WorkSpace\src\golib\uchystrixhttp\quanshi.cer`, `D:\Go\WorkSpace\src\golib\uchystrixhttp\quanshi.key`)
	if err != nil {
		t.Error("Loadx509keypair err:", err)
		return
	}
	_, b, err := SSLRequest("GET", "https://testcloudb.quanshi.com/bee/client/download-html/download.php", nil, cliCrt)
	t.Log(string(b))
	t.Error(err)
}

func TestRequest7(t *testing.T) {
	InitTest()
	cliCrt, err := tls.LoadX509KeyPair(`D:\Go\WorkSpace\src\golib\uchystrixhttp\quanshi.cer`, `D:\Go\WorkSpace\src\golib\uchystrixhttp\quanshi.key`)
	if err != nil {
		t.Error("Loadx509keypair err:", err)
		return
	}
	_, b, err := SSLRequestTimeout("GET", "https://testcloudb.quanshi.com/bee/client/download-html/download.php", nil, time.Duration(1000)*time.Millisecond, time.Duration(1000)*time.Millisecond, cliCrt)
	t.Log(string(b))
	t.Error(err)
}

func TestRequest8(t *testing.T) {
	var err error
	InitTest()
	statusCode, b, err := RequestForm("POST", "http://family.quanshi.com/new_index.php?m=Index&c=Public&a=login", strings.NewReader("username=zheng.guan&password=Gz1985@10@14"))
	t.Log(string(b))
	t.Log(statusCode)
	t.Error(err)
}

func TestRequest9(t *testing.T) {
	var err error
	InitTest()
	uchystrixhttp.HystrixEnable = true
	uchystrixhttp.HystrixStatEnable = true

	uchystrixhttp.HystrixConfig = `{"ucc_command": {"timeout": 1000,"max_concurrent_requests": 20,"request_volume_threshold": 20,"sleep_window": 5000,"error_percent_threshold": 50},"ums_command": {"timeout": 1000,"max_concurrent_requests": 20,"request_volume_threshold": 20,"sleep_window": 5000,"error_percent_threshold": 50}}`
	uchystrixhttp.HystrixCommandConfigString = `{"baidu":"ucc_command","user/check":"ucc_command","rs/users/id/in":"ums_command"}`
	uchystrixhttp.Init()
	statusCode, b, err := RequestXML("GET", "https://www.baidu.com/img/baidu.svg", nil)
	t.Log(string(b))
	t.Log(statusCode)
	t.Error(err)
}

func TestRequest10(t *testing.T) {
	var err error
	uchystrixhttp.HystrixEnable = true
	uchystrixhttp.HystrixStatEnable = true

	uchystrixhttp.HystrixConfig = `{"ucc_command": {"timeout": 100,"max_concurrent_requests": 20,"request_volume_threshold": 20,"sleep_window": 5000,"error_percent_threshold": 50},"ums_command": {"timeout": 1000,"max_concurrent_requests": 20,"request_volume_threshold": 20,"sleep_window": 5000,"error_percent_threshold": 50}}`
	uchystrixhttp.HystrixCommandConfigString = `{"baidu":"ucc_command","user/check":"ucc_command","rs/users/id/in":"ums_command"}`
	uchystrixhttp.Init()
	InitTest()
	statusCode, b, err := RequestTimeoutWithHeader("GET", "https://www.baidu.com/img/baidu.svg", nil, time.Duration(1000)*time.Millisecond, time.Duration(1000)*time.Millisecond, map[string]string{"Content-Type": "application/xml"})
	t.Log(string(b))
	t.Log(statusCode)

	t.Error(err)
}
