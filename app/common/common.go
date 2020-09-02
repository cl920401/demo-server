package common

import "github.com/gin-gonic/gin"

const (
	HEADER_LOGIN_TOKEN = "X-Login-Token"
	HEADER_USER_UUID   = "User-UUID"
	HEADER_REQUEST_ID  = "Request-ID"
)

const (
	CONTEXT_VALUE_LOGIN_TOKEN = "login_token"
	CONTEXT_VALUE_UUID        = "uuid"
	CONTEXT_VALUE_REQUEST_ID  = "request_id"
	CONTEXT_VALUE_REQUEST_ERR = "request_error"
)

func GetContextString(c *gin.Context, key string) string {
	if c == nil {
		return ""
	}
	value, has := c.Get(key)
	if !has {
		return ""
	}
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return ""
}
