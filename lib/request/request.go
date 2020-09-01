package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"demo-server/utils/tools"
	"net/http"
	"strings"
	"time"

	"demo-server/lib/log"

	"github.com/parnurzeal/gorequest"
)

var (
	HttpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

// Request : 请求重试配置
type Request struct {
	TimeOut   time.Duration
	RetryNum  int
	RetryTime time.Duration
}

func (r *Request) request(api string, method string, params interface{}, header map[string]string) ([]byte, error) {
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
			// 默认使用json, 包含Content-Type: application/x-www-form-urlencoded时转换form参数
			if header["Content-Type"] == gorequest.Types[gorequest.TypeForm] {
				request.Type(gorequest.TypeForm).Send(params)
			} else {
				request.Send(params)
			}
		}
	}

	//设置请求头
	for hk, hv := range header {
		request.Set(hk, hv)
	}

	//发送请求
	ret, body, errs := request.EndBytes()
	log.Debug(request.AsCurlCommand())
	log.Debug(tools.Bytes2str(body))
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

// Request : 发送http请求 返回结果解析到json中 注意res的类型应该是指针
func (r *Request) Request(api string, method string, params interface{}, header map[string]string, res interface{}) error {
	body, err := r.request(api, method, params, header)
	if err != nil {
		return fmt.Errorf("request error:[%w]", err)
	}
	return json.Unmarshal(body, res)
}

// Request : 发送http请求 返回结果解析到json中 注意res的类型应该是指针
func (r *Request) RequestWithFile(api string, method string, params map[string]string, file io.Reader, fileName, fieldName string, res interface{}) error {
	body, err := r.requestWithFile(api, method, params, file, fileName, fieldName)
	if err != nil {
		return fmt.Errorf("request error:[%w]", err)
	}
	return json.Unmarshal(body, res)
}

func (r *Request) requestWithFile(api string, method string, params map[string]string, file io.Reader, fileName, fieldName string) ([]byte, error) {
	//body := new(bytes.Buffer)

	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	for key, val := range params {
		_ = body_writer.WriteField(key, val)
	}

	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	//close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	formFile, err := body_writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	err = body_writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, api, body_buf)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Content-Type", writer.FormDataContentType())

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	log.Debug(req.URL)
	log.Debug(tools.Bytes2str(content))
	if err != nil {
		return nil, err
	}
	return content, nil
}
