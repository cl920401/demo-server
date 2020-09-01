package router

import (
	"demo-server/app/example/controller/test"
	"demo-server/app/internal/context"
	"demo-server/app/internal/router"

	"github.com/gin-gonic/gin"

	//"demo-server/def/router"
	"net/http"
)

var URL = map[string][]router.Info{
	// ################### 测试 ##################
	"/robot/v1/test": {
		{
			Method:  http.MethodGet,
			Name:    "测试",
			Handler: test.Test,
		},
	},
}

//设置路由
func API(r *gin.Engine) {
	for path, act := range URL {
		for _, v := range act {
			r.Handle(v.Method, path, context.Handle(v.Handler))
		}
	}
}
