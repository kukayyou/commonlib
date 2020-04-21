package ucutil

import (
	"fmt"
	"gnetis.com/golang/core/golib/uchttp"
	"gnetis.com/golang/core/golib/uclog"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

/**
* 描述：curl command
* param：method GET/POST, USE GET WHEN CALL AUTH2
* requrl: url
* reqstr: key=val&key1=val1, only access_token=val when call AUTH2
* return：int return code
* return：error 错误
*
 */

func HttpCurl(method string, requrl string, reqstr string, send_format string, accept_format string) (int, []byte, error) {

	request, err := http.NewRequest(method, requrl, strings.NewReader(reqstr))

	if err != nil {
		httpErr := fmt.Errorf("construct http request failed, requrl = %s, err:%s", requrl, err.Error())
		uclog.Critical(httpErr.Error())
		return -1, nil, httpErr
	}
	if send_format != "" {
		request.Header.Add("Content-Type", send_format)
	}
	if accept_format != "" {
		request.Header.Add("Accept", accept_format)
	}

	if strings.Contains(reqstr, "access_token=") {
		token := fmt.Sprintf("Bearer %s", strings.Replace(reqstr, "access_token=", "", -1))
		uclog.Debug("token=%s", token)
		request.Header.Add("Authorization", token)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, uchttp.HTTP_TIMEOUT_CONNECT) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(uchttp.HTTP_TIMEOUT_DEADLINE)) //设置发送接收数据超时
				return c, nil
			},
			DisableKeepAlives: true,
		},
	}

	response, err := client.Do(request)

	if err != nil {
		httpErr := fmt.Errorf("http request failed, url: %s, error:%s", requrl, err.Error())
		uclog.Critical(httpErr.Error())
		return -1, nil, httpErr
	}

	if response != nil {

		//uclog.Info("curl request string is %s, url is %s, response code is %d\n", reqstr, requrl, response.StatusCode)

		defer response.Body.Close()

		respBody, _ := ioutil.ReadAll(response.Body)

		//uclog.Info("curl reponse is: \n%s", string(respBody))

		if response.StatusCode > 300 {
			httpErr := fmt.Errorf("http request failed, url: %s, error:%s", requrl, err.Error())
			uclog.Critical(httpErr.Error())
			return response.StatusCode, respBody, httpErr
		}

		return response.StatusCode, respBody, nil
	}

	return -1, nil, fmt.Errorf("%s to %s, empty response", method, requrl)
}
