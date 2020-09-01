package kong

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"

	"demo-server/lib/util"
	"net/url"
	"strings"
	"time"
)

const (
	Algorithm_SHA1   = "hmac-sha1"
	Algorithm_SHA256 = "hmac-sha256"
	Algorithm_SHA384 = "hmac-sha384"
	Algorithm_SHA512 = "hmac-sha512"

	AUTH_HEADER_NAME        = "Authorization"
	AUTH_HEADER_VALUE_FMT   = `hmac username="%s", algorithm="%s", headers="%s", signature="%s"`
	BODY_DIGEST_HEADER_NAME = "Digest"
	DATE_HEADER_NAME        = "X-Date"
)

type HMACAuth struct {
	Config HMACAuthConfig
}

type HMACAuthConfig struct {
	AppKey         string
	AppSecret      string
	EnforceHeaders []string
	Algorithm      string

	TimeOut   time.Duration
	RetryNum  int
	RetryTime time.Duration
}

func NewHMACAuth(conf HMACAuthConfig) *HMACAuth {
	cli := &HMACAuth{
		Config: conf,
	}
	return cli
}

func (this *HMACAuth) Get(api string, params interface{}, header map[string]string, jsonRes interface{}) ([]byte, error) {
	return this.doRequest(api, http.MethodGet, params, header, jsonRes)
}

func (this *HMACAuth) PostFormUrlencoded(api string, params url.Values, header map[string]string, jsonRes interface{}) ([]byte, error) {
	bts := []byte(params.Encode())
	header["Content-Type"] = "application/x-www-form-urlencoded"
	return this.doRequest(api, http.MethodPost, bts, header, jsonRes)
}

// PostFormData需要支持文件上传，暂时没有通过网关上传文件的需求
//func (this *HMACAuth) PostFormData(api string, params url.Values, header map[string]string, jsonRes interface{}) ([]byte, error) {
//	header["Content-Type"] = "multipart/form-data; boundary=----WebKitFormBoundaryrGKCBY7qhFd3TrwA"
//	return nil, nil
//}

func (this *HMACAuth) PostJSON(api string, params interface{}, header map[string]string, jsonRes interface{}) ([]byte, error) {
	bts, _ := json.Marshal(params)
	header["Content-Type"] = "application/json"
	return this.doRequest(api, http.MethodPost, bts, header, jsonRes)
}

// RequestBindJSON 发送http请求 返回结果解析到json中
func (this *HMACAuth) doRequest(api string, method string, params interface{}, header map[string]string, jsonRes interface{}) (body []byte, err error) {
	var requestBody []byte
	if strings.ToUpper(method) == http.MethodPost {
		requestBody = params.([]byte) // 私有函数内部调用，此处不必校验
		params = util.Bytes2str(requestBody)
	}

	// 获取签名后的authHeader
	var authHeader map[string]string
	authHeader, err = this.signature(requestBody, header)
	if err != nil {
		return
	}

	// 复用已有的http请求组件
	r := util.Request{
		TimeOut:           this.Config.TimeOut,
		RetryNum:          this.Config.RetryNum,
		RetryTime:         this.Config.RetryTime,
		BounceToRawString: true,
	}
	body, err = r.Request(api, method, params, authHeader)
	if err != nil {
		return
	}

	// 如果有传入jsonRes，将执行json解析到jsonRes第一个变量中
	// 由于外部调用需要区分传输方式，如果分别实现BindJSON需要实现的方法太多，此处合并了r.BindJSON的功能
	if jsonRes != nil {
		err = json.Unmarshal(body, &jsonRes)
		if err != nil {
			return
		}
	}
	return
}

func (this *HMACAuth) signature(requestBody []byte, header map[string]string) (authHeader map[string]string, err error) {
	authHeader = make(map[string]string)
	headerKeyList := []interface{}{}
	for k, v := range header {
		authHeader[k] = v
	}
	// Digest body不为空则参与计算
	if requestBody != nil && len(requestBody) > 0 {
		has := sha256.New()
		has.Write([]byte(requestBody))
		bs := has.Sum(nil)
		authHeader[BODY_DIGEST_HEADER_NAME] = "SHA-256=" + util.Base64EncodeToString(bs)
		headerKeyList = append(headerKeyList, BODY_DIGEST_HEADER_NAME)
	}

	_, hasXData := header["X-Date"]
	_, hasData := header["Date"]

	// 传入header中X-Date、Date均未包含时，补全X-Date
	if !hasData && !hasXData {
		authHeader[DATE_HEADER_NAME] = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		headerKeyList = append(headerKeyList, DATE_HEADER_NAME)
	}

	//  将配置中要求校验的Headers与默认headers合并，同时去重
	for _, v := range this.Config.EnforceHeaders {
		headerKeyList = append(headerKeyList, v)
	}
	headerKeyList = util.SliceUnique(headerKeyList)

	// 合并待签名字符串
	// 有些签名文档要求ASCII排序，其实不必
	// hmac签名只要求`hmac username="%s", algorithm="%s", headers="%s", signature="%s"`中的headers顺序和签名字串key顺序一致
	signatureStr := ""
	headersStr := ""
	for _, item := range headerKeyList {
		key := item.(string) // 私有函数内部调用，此处不必校验
		if val, ok := authHeader[key]; ok {
			signatureStr += key + ": " + val + "\n"
			headersStr += key + " "
		} else {
			err = fmt.Errorf("No header %s defined", key)
			return
		}
	}
	signatureStr = signatureStr[:len(signatureStr)-1]
	headersStr = headersStr[:len(headersStr)-1]

	// 通过客户端设定的签名算法，选择调用方法
	var shaNewFunc func() hash.Hash
	switch this.Config.Algorithm {
	case Algorithm_SHA1:
		shaNewFunc = sha1.New
	case Algorithm_SHA256:
		shaNewFunc = sha256.New
	case Algorithm_SHA384:
		shaNewFunc = sha512.New384
	case Algorithm_SHA512:
		shaNewFunc = sha512.New
	default:
		err = fmt.Errorf("Unknown algorithm")
		return
	}

	//完成签名计算
	mac := hmac.New(shaNewFunc, []byte(this.Config.AppSecret))
	mac.Write([]byte(signatureStr))
	bs := mac.Sum(nil)

	// 名字符串
	signature := util.Base64EncodeToString(bs)

	// 鉴权token
	Auth := fmt.Sprintf(AUTH_HEADER_VALUE_FMT,
		this.Config.AppKey,
		this.Config.Algorithm,
		headersStr,
		signature)
	authHeader[AUTH_HEADER_NAME] = Auth
	return
}
