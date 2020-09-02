package middleware

import (
	"github.com/gin-gonic/gin"
)

// Hook : 请求中间件，所有请求会先打到这个方法
func Hook() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
