package costool

import (
	"io"
	"os"
	"strings"
	"testing"

	"demo-server/lib/config"
	"demo-server/lib/log"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func TestMain(m *testing.M) {

	if err := config.LoadFile("../../configs/user/user_dev_conf.json"); err != nil {

		_ = config.LoadEnv("ETCD_")

		err := config.LoadEtcd(
			config.Get("etcd.prefix").String("/root/nlp/demo/user/dev"),
			config.Get("etcd.user").String("nlp"),
			config.Get("etcd.pwd").String("nlp"),
			strings.Split(config.Get("etcd.endpoints").String("172.16.101.128:2379,172.16.101.100:2379,172.16.101.78:2379"), ",")...,
		)
		log.Error("conf.LoadFile() error:", err)
	}
	InitCos()
	os.Exit(m.Run())
}

func TestNewCosClient(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test NewCosClient"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCosClient(); got == nil {
				t.Errorf("NewCosClient() = %v", got)
			}
		})
	}
}

func TestCosUpload(t *testing.T) {
	type args struct {
		path   string
		name   string
		f      io.Reader
		client *cos.Client
	}
	file, err := os.Open("./test.txt")
	if err != nil {
		t.Errorf("ReadFile error: %s", err.Error())
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test CosUpload", args{path: "map/test/", name: "test.txt", f: file, client: NewCosClient()}, "map/test/test.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CosUpload(tt.args.path, tt.args.name, tt.args.f, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("CosUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CosUpload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDownloadUrl(t *testing.T) {
	type args struct {
		name     string
		filename string
		client   *cos.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test GetDownloadUrl", args{"map/test/test.txt", "test.txt", NewCosClient()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDownloadUrl(tt.args.name, tt.args.filename, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDownloadUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetDownloadUrl() = %v", got)
		})
	}
}

func TestGetUrl(t *testing.T) {
	type args struct {
		name   string
		client *cos.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test GetUrl", args{"map/test/test.txt", NewCosClient()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUrl(tt.args.name, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetUrl() = %v", got)
		})
	}
}

func TestDownload(t *testing.T) {
	type args struct {
		name   string
		client *cos.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test Download", args{"map/test/test.txt", NewCosClient()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Download(tt.args.name, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Download() = %v", got)
		})
	}
}

func TestDownloadFile(t *testing.T) {
	type args struct {
		name   string
		client *cos.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test DownloadFile", args{"map/test/test.txt", NewCosClient()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadFile(tt.args.name, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("DownloadFile() = %v", got)
		})
	}
}

func TestCosUploadPrivate(t *testing.T) {
	type args struct {
		path   string
		name   string
		f      io.Reader
		client *cos.Client
	}
	file, err := os.Open("./test.txt")
	if err != nil {
		t.Errorf("ReadFile error: %s", err.Error())
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"TestCosUploadPrivate", args{path: "map/test/", name: "test.txt", f: file, client: NewCosClient()}, "map/test/test.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CosUploadPrivate(tt.args.path, tt.args.name, tt.args.f, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("CosUploadPrivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CosUploadPrivate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosUploadEx(t *testing.T) {
	type args struct {
		path   string
		name   string
		f      io.Reader
		client *cos.Client
		opt    *cos.ObjectPutOptions
	}
	file, err := os.Open("./test.txt")
	if err != nil {
		t.Errorf("ReadFile error: %s", err.Error())
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test CosUploadEx", args{path: "map/test/", name: "test.txt", f: file, client: NewCosClient()}, "map/test/test.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CosUploadEx(tt.args.path, tt.args.name, tt.args.f, tt.args.client, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("CosUploadEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CosUploadEx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
