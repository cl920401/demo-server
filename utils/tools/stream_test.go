package tools

import (
	"strconv"
	"testing"
	"time"
)

func TestGenerateGuid(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"TestGenerateGuid", "fe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateGUID()
			t.Log(got)
		})
	}
}

func TestEncrypt(t *testing.T) {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	type args struct {
		data []byte
		key  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestEncrypt",
			args: args{
				data: []byte(strconv.FormatInt(time.Now().In(cstSh).Unix(), 10)),
				key:  "demo-reporter",
			},
			want: "lpecZWtqbGdrYw==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encrypt(tt.args.data, tt.args.key)
			t.Log(got)
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		data string
		key  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDecrypt",
			args: args{
				data: "FhccZWtqbGdrYw==",
				key:  "aXnpscNiWvVfYmbd",
			},
			want: "1583385190",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decrypt(tt.args.data, tt.args.key)
			t.Log(got)
		})
	}
}
