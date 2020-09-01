package feishu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

//飞书通知发送专用 Request 避免循环import
type Request struct {
	TimeOut   time.Duration
	RetryNum  int
	RetryTime time.Duration
}

// 发送http请求
// 重试三次间隔一分钟
// api地址 请求方法 请求参数 请求头设置 返回数据解析
func (r *Request) Request(api string, method string, params interface{}, header map[string]string) ([]byte, error) {
	if api == "" {
		return nil, errors.New("api url null")
	}
	//创建request请求对象
	//设置重试状态
	request := gorequest.New().CustomMethod(strings.ToUpper(method), api)
	if r.TimeOut > 0 {
		request = request.Timeout(r.TimeOut)
	}
	if r.RetryNum > 0 {
		request = request.Retry(r.RetryNum, r.RetryTime, http.StatusBadRequest, http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout)
	}
	//设置请求参数
	if params != nil {
		if strings.ToUpper(method) == http.MethodGet {
			request.Query(params)
		} else {
			request.Send(params)
		}
	}
	//设置请求头
	for hk, hv := range header {
		request.Set(hk, hv)
	}
	//发送请求
	ret, body, errs := request.EndBytes()
	//验证错误
	if len(errs) > 0 {
		return nil, errs[0]
	}
	//验证返回状态
	if ret.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error %d", ret.StatusCode)
	}
	if len(body) == 0 {
		return nil, errors.New("body null")
	}
	return body, nil
}

// RequestBindJSON 发送http请求 返回结果解析到json中
func (r *Request) BindJSON(api string, method string, params interface{}, header map[string]string, res interface{}) error {
	body, err := r.Request(api, method, params, header)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &res)
}

//Download 下载文件
func (r *Request) Download(fileURL string, method string, params interface{}, header map[string]string, dst string) error {
	if fileURL == "" {
		return errors.New("fileURL url null")
	}
	//创建request请求对象
	//设置重试状态
	request := gorequest.New().CustomMethod(strings.ToUpper(method), fileURL)
	if r.TimeOut > 0 {
		request = request.Timeout(r.TimeOut)
	}
	if r.RetryNum > 0 {
		request = request.Retry(r.RetryNum, r.RetryTime, http.StatusBadRequest, http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout)
	}
	//设置请求参数
	if params != nil {
		if strings.ToUpper(method) == http.MethodGet {
			request.Query(params)
		} else {
			request.Send(params)
		}
	}
	//设置请求头
	for hk, hv := range header {
		request.Set(hk, hv)
	}
	//发送请求
	req, err := request.MakeRequest()
	if err != nil {
		return err
	}

	ret, err := request.Client.Do(req)
	if err != nil {
		return err
	}
	//验证返回状态
	if ret.StatusCode != http.StatusOK {
		return fmt.Errorf("http code error %d", ret.StatusCode)
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	defer ret.Body.Close()
	_, err = io.Copy(out, ret.Body)
	return err
}
