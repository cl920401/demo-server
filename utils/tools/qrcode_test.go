package tools

import "testing"

func TestQRFile(t *testing.T) {
	type args struct {
		url  string
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestQRFile",
			args: args{
				url:  "https://minicdn.ainirobot.com/web/mini_app/download.html?bind_code=58c598078caac323f122d2ad678f1a34&robot_sn=KTS17P480336",
				file: "/tmp/test.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := QRFile(tt.args.url, tt.args.file, 256); (err != nil) != tt.wantErr {
				t.Errorf("QRFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
