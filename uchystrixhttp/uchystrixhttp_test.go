package uchystrixhttp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gnetis.com/golang/core/golib/uclog"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func InitTest() {
	uclog.Initialize("log", `D:\Go\WorkSpace\src\golib\uchystrixhttp\a.log`, "debug")
}

func TestRequest1(t *testing.T) {
	InitTest()
	_, b, _, err := New("GET", "http://www.baidu.com", nil).SetHeader(map[string]string{"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"}).Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest2(t *testing.T) {
	InitTest()
	_, b, _, err := New("GET", "http://www.baidu.com", nil).Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest3(t *testing.T) {
	InitTest()
	_, b, _, err := New("GET", "http://www.baidu.com", nil).SetTimeout(time.Duration(100)*time.Millisecond, time.Duration(100)*time.Millisecond).Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest4(t *testing.T) {
	InitTest()
	_, b, _, err := New("GET", "http://family.quanshi.com/new_index.php?m=&&m=Hr&c=index&a=index", nil).SetCookies([]*http.Cookie{&http.Cookie{Name: "PHPSESSID", Value: "79d2rliq62lkv3cqeerda4pok2"}}).Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest5(t *testing.T) {
	InitTest()
	_, b, _, err := New("GET", "http://127.0.0.1:3388/", nil).SetAuth("upupw", "upupw").Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest6(t *testing.T) {
	InitTest()
	_, b, _, err := New("GET", "http://www.google.com/", nil).SetProxy("http://127.0.0.1:1080/pac?t=20180105113011554&secret=4p8ERMHkdh27WvucpSdWS0wIwai43Fpsc+U7Pp/3T/A=").Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest8(t *testing.T) {
	InitTest()
	cliCrt, err := tls.LoadX509KeyPair(`D:\Go\WorkSpace\src\golib\uchystrixhttp\quanshi.cer`, `D:\Go\WorkSpace\src\golib\uchystrixhttp\quanshi.key`)
	if err != nil {
		t.Error("Loadx509keypair err:", err)
		return
	}
	req := New("GET", "https://testcloudb.quanshi.com/bee/client/download-html/download.php", nil)
	req.SetCerts([]tls.Certificate{cliCrt})
	_, b, _, err := req.Request(func(e error) (int, []byte, *http.Response, error) {
		return 403, nil, nil, fmt.Errorf("request error:%s", e.Error())
	})
	t.Log(string(b))
	t.Error(err)
}

func TestRequest9(t *testing.T) {
	InitTest()
	HystrixEnable = true
	HystrixStatEnable = true

	HystrixConfig = `{"ucc_command": {"timeout": 60,"max_concurrent_requests": 20,"request_volume_threshold": 20,"sleep_window": 5000,"error_percent_threshold": 50},"ums_command": {"timeout": 1000,"max_concurrent_requests": 20,"request_volume_threshold": 20,"sleep_window": 5000,"error_percent_threshold": 50}}`
	HystrixCommandConfigString = `{"message/msgsend":"ucc_command","user/check":"ucc_command","rs/users/id/in":"ums_command"}`
	Init()
	for k, v := range hystrix.GetCircuitSettings() {
		t.Error(k, v)
	}
}

var data []map[string]string
var data2 map[string]string

func TestRequest10(t *testing.T) {
	HystrixCommandConfigString := `[{"1generalAuthentication":"login"},{"2ums-oauth2":"oauth"},{"3umsapi":"ums"},{"4solr":"search"},{"5uccapi":"ucc"},{"6presenceapi":"presence"}] `

	json.Unmarshal([]byte(HystrixCommandConfigString), &data)
	log.Printf("%v", data)
	for _, v := range data {
		for tag, name := range v {
			log.Printf("%s  ->   %s", tag, name)
		}
	}
	t.Error(data)

}

func TestRequest11(t *testing.T) {
	HystrixCommandConfigString := `{"1generalAuthentication":"login","2ums-oauth2":"oauth","3umsapi":"ums","4solr":"search","5uccapi":"ucc","6presenceapi":"presence"}`
	// HystrixCommandConfigString := `[{"1generalAuthentication":"login"},{"2ums-oauth2":"oauth","3umsapi":"ums"},{"4solr":"search"},{"5uccapi":"ucc"},{"6presenceapi":"presence"}] `
	if err := json.Unmarshal([]byte(HystrixCommandConfigString), &data); err != nil {
		if err := json.Unmarshal([]byte(HystrixCommandConfigString), &data2); err == nil {
			data = append(data, data2)
		}
	}
	var commandName string
	for _, v := range data {
		for tag, name := range v {
			log.Printf("%s  ->   %s", tag, name)

			if strings.Contains("http://fdsfsklf.dd/3umsapi/1generalAuthentication", tag) {
				commandName = name
				break
			}
			// if tag != "" {
			// 	commandName = name
			// 	break
			// }
		}
		if commandName != "" {
			break
		}
	}
	log.Println(commandName)
	t.Error(1)

}
