package redis

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"demo-server/lib/mysql"

	"demo-server/lib/config"
	"demo-server/lib/log"
)

func TestMain(m *testing.M) {

	err := config.LoadFile("../../configs/gateway/gateway_dev_conf.json")
	if err != nil {
		_ = config.LoadEnv("ETCD_")

		err := config.LoadEtcd(
			config.Get("etcd.prefix").String("/root/nlp/demo/user/dev"),
			config.Get("etcd.user").String("nlp"),
			config.Get("etcd.pwd").String("nlp"),
			strings.Split(config.Get("etcd.endpoints").String("172.16.101.128:2379,172.16.101.100:2379,172.16.101.78:2379"), ",")...,
		)
		log.Error("conf.LoadFile err:", err)
	}
	InitRedis()
	mysql.InitDB()
	os.Exit(m.Run())
}

func TestConnect_HGetKeyString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"TestConnect_HGetKeyString", args{KeyPrefix}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisConn := Default()
			got, err := redisConn.GetKeys(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGetKeyString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
