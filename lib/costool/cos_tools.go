package costool

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"demo-server/lib/config"
	"demo-server/lib/log"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var setting struct {
	CosHost   string `json:"cos_host"`
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

var isInit bool

//InitCos cos初始化
func InitCos() {
	if isInit {
		return
	}
	conf := config.Get("cos")
	if err := conf.Scan(&setting); err != nil {
		log.Fatal("Error parsing cos configuration file ", err)
		return
	}
	isInit = true
}

func NewCosClient() *cos.Client {
	u, _ := url.Parse(setting.CosHost)
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  setting.SecretID,
			SecretKey: setting.SecretKey,
		},
	})
}

// 上传一个私有读文件
func CosUploadPrivate(path string, name string, f io.Reader, client *cos.Client) (string, error) {
	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	_, err := client.Object.Put(context.Background(), path+name, f, opt)
	if err != nil {
		return "", err
	}
	return path + name, nil
}

func CosUploadEx(path string, name string, f io.Reader, client *cos.Client, opt *cos.ObjectPutOptions) (string, error) {
	_, err := client.Object.Put(context.Background(), path+name, f, opt)
	if err != nil {
		return "", err
	}
	return path + name, nil
}

func CosUploadPng(filePath string, r io.Reader, client *cos.Client) error {
	if _, err := client.Object.Put(context.Background(), filePath, r, nil); err != nil {
		return err
	}
	return nil
}

func CosUpload(path string, name string, r io.Reader, client *cos.Client) (string, error) {
	_, err := client.Object.Put(context.Background(), path+name, r, nil)
	if err != nil {
		return "", err
	}
	return path + name, nil
}

func GetDownloadUrl(name string, filename string, client *cos.Client) (string, error) {
	opt := &cos.ObjectGetOptions{
		ResponseContentType:        "application/octet-stream",
		ResponseContentDisposition: fmt.Sprintf("attachment;filename=%s", filename),
	}
	presignedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodGet, name, setting.SecretID, setting.SecretKey, time.Hour, opt)
	if err != nil || presignedURL == nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func GetUrl(name string, client *cos.Client) (string, error) {
	presignedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodGet, name, setting.SecretID, setting.SecretKey, time.Hour, nil)
	if err != nil || presignedURL == nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func PutUrl(name string, client *cos.Client) (string, error) {
	presignedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodPut, name, setting.SecretID, setting.SecretKey, time.Hour, nil)
	if err != nil || presignedURL == nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func Download(name string, client *cos.Client) ([]byte, error) {
	resp, err := client.Object.Get(context.Background(), name, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	return ioutil.ReadAll(resp.Body)
}

func DownloadFile(name string, client *cos.Client) (string, error) {
	filePath := "/tmp/cos/" + name

	if err := os.MkdirAll(path.Dir(filePath), 0777); err != nil {
		log.Error("Error creating directory:", err)
		return "", err
	}

	resp, err := client.Object.GetToFile(context.Background(), name, filePath, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return filePath, nil
}
