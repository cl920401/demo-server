package test

import (
	"demo-server/app/internal/context"
	"demo-server/app/internal/message"
	"demo-server/lib/costool"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	conf "demo-server/lib/config"
	"demo-server/lib/log"
	"demo-server/lib/mysql"
	"demo-server/lib/redis"
)

func TestMain(m *testing.M) {
	err := conf.LoadEtcd(
		conf.Get("etcd.prefix").String(""),
		conf.Get("etcd.user").String(""),
		conf.Get("etcd.pwd").String(""),
		strings.Split(conf.Get("etcd.endpoints").String(""), ",")...,
	)

	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	redis.InitRedis()
	mysql.InitDB()
	costool.InitCos()
	os.Exit(m.Run())
}

func TestTest(t *testing.T) {
	testCase := []struct {
		req *http.Request
		res message.ResponseData
	}{
		{
			req: httptest.NewRequest("GET", "/test?Arg1=111", nil),
			res: message.ResponseData{Header: message.ResponseHeader{Code: message.CodeOK}},
		},
		{
			req: httptest.NewRequest("GET", "/test?Arg1=aaaa", nil),
			res: message.ResponseData{Header: message.ResponseHeader{Code: message.CodeUnknowErr}},
		},
		{
			req: httptest.NewRequest("GET", "/test", nil),
			res: message.ResponseData{Header: message.ResponseHeader{Code: message.CodeParamErr}},
		},
	}

	for _, c := range testCase {
		ctx := context.CreateTestContext(c.req)
		resBody := Test(ctx)
		bts, _ := json.Marshal(resBody)
		if resBody.Header.Code != c.res.Header.Code {
			t.Logf("Not expected resBody=%s c.res=%+v", bts, c.res)
			t.Fail()
		} else {
			t.Logf("passed resBody=%s", bts)
		}
	}
}
