package context

import (
	"github.com/gin-gonic/gin"
)

const (
	HEADER_LOGIN_TOKEN = "X-Login-Token"
	HEADER_OV_TOKEN    = "X-OV-Token"
	HEADER_DEVICE_TYPE = "Device-Type"
	HEADER_DEVICE_ID   = "Device-ID"
	HEADER_REQUEST_ID  = "Request-ID"
	HEADER_USER_UUID   = "User-UUID"
	HEADER_FAMILY_ID   = "Family-ID"
)

const (
	CONTEXT_VALUE_LOGIN_TOKEN = "login_token"
	CONTEXT_VALUE_OV_TOKEN    = "ov_token"
	CONTEXT_VALUE_DEVICE_TYPE = "device_type"
	CONTEXT_VALUE_DEVICE_ID   = "device_id"
	CONTEXT_VALUE_REQUEST_ID  = "request_id"
	CONTEXT_VALUE_REQUEST_ERR = "request_error"
	CONTEXT_VALUE_UUID        = "uuid"
	CONTEXT_VALUE_FAMILY_ID   = "family_id"
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
