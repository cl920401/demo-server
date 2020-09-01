module demo-server

go 1.13

require (
	github.com/Masterminds/squirrel v1.4.0
	github.com/Shopify/sarama v1.27.0
	github.com/aws/aws-sdk-go v1.34.14
	github.com/elazarl/goproxy v0.0.0-20200809112317-0581fc3aee2d // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9 // indirect
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/facebookgo/httpdown v0.0.0-20180706035922-5979d39b15c2 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/stats v0.0.0-20151006221625-1b76add642e4 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/hashicorp/go-uuid v1.0.2
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/olivere/elastic v6.2.35+incompatible
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tencentyun/cos-go-sdk-v5 v0.7.8
	github.com/typa01/go-utils v0.0.0-20181126045345-a86b05b01c1e
	github.com/urfave/cli v1.22.4
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/zap v1.15.0
	google.golang.org/grpc v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	moul.io/http2curl v1.0.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
