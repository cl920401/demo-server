package request

import (
	"testing"
	"time"
)

//{
//	"sid": "2222fbfd-aa30-4013-b4aa-b98733b1",
//	"asr": {
//		"word": "今天天气怎么样"
//	},
//	"semantic": [{
//		"english_domain": "weather",
//		"domain": "天气",
//		"intent": "get_weather",
//		"slots": {
//			"time": [{
//				"slot_type": "TIME",
//				"text": "今天",
//				"value": {
//					"begin": {
//						"date": {
//							"day": 29,
//							"month": 1,
//							"year": 2018
//						}
//					},
//					"sub_type": 0,
//					"type": 1
//				}
//			}]
//		},
//		"source": "",
//		"semantics_flag": 1
//	}],
//	"context": {
//		"client_id": "orion.ovs.client.1514259512471",
//		"enterprise_id": "orion.ovs.entprise.4520242975",
//		"group_id": "ovs.group.156015364495017",
//		"model": "CM-GB01L",
//		"dt": "1517221791000",
//		"deviceid": "L4A0D011719A38NCCA",
//		"union_access_token": "",
//		"device_type": "1",
//		"os_type": "android",
//		"lat": "",
//		"lng": ""
//	}
//}

func TestRequest_BindJSON(t *testing.T) {
	type fields struct {
		TimeOut   time.Duration
		RetryNum  int
		RetryTime time.Duration
	}
	type args struct {
		api    string
		method string
		params interface{}
		header map[string]string
		res    interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "TestRequest_BindJSON",
			fields: fields{},
			args: args{
				api:    "http://ovstest.ainirobot.com:8091/dialogue",
				method: "POST",
				params: nil,
				header: nil,
				res:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				TimeOut:   tt.fields.TimeOut,
				RetryNum:  tt.fields.RetryNum,
				RetryTime: tt.fields.RetryTime,
			}
			if err := r.Request(tt.args.api, tt.args.method, tt.args.params, tt.args.header, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
