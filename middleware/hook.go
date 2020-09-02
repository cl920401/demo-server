package middleware

import (
	"demo-server/app/common"
	"demo-server/lib/log"
	"demo-server/utils/tools"

	"github.com/gin-gonic/gin"
)

// Hook : 请求中间件，所有请求会先打到这个方法
func Hook() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置header变量
		requestID := c.GetHeader(common.HEADER_REQUEST_ID)
		if requestID == "" {
			requestID = tools.GetRandomString(3)
		}
		c.Set(common.CONTEXT_VALUE_LOGIN_TOKEN, c.GetHeader(common.HEADER_LOGIN_TOKEN))
		c.Set(common.CONTEXT_VALUE_UUID, c.GetHeader(common.HEADER_USER_UUID))
		c.Set(common.CONTEXT_VALUE_REQUEST_ID, requestID)
		c.Next()
		// 输出接口错误日志 需要在controller层写入error到CONTEXT_VALUE_REQUEST_ERR
		if err, has := c.Get(common.CONTEXT_VALUE_REQUEST_ERR); has && err != nil {
			log.Warn(" error request id:", requestID, " error info:", err)
		}
	}
}
