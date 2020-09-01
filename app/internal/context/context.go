package context

import (
	"fmt"
	"demo-server/app/internal/message"
	"demo-server/utils/tools"
	"net/http"
	"net/http/httptest"
	"runtime"

	"demo-server/lib/log"
	"github.com/gin-gonic/gin"
)

const (
	DEVICE_TYPE_IOS     = "1"
	DEVICE_TYPE_ANDROID = "2"
	DEVICE_TYPE_ROBOT   = "3"
)

//自定义HandlerFunc
type HandlerFunc func(c *Context) *message.ResponseData

type IContextAPI interface {
	GetRequestId() string
	GetUserID() string
	GetFamilyID() string
	GetDeviceID() string
	GetDeviceType() string
	IsIOS() bool
	IsAndroid() bool
	IsRobot() bool
	IsAPP() bool
	//GetRequestError() error
	//SetRequestError(err error)
	LogFormat(format string, a ...interface{}) string
}

type IContextBase interface {
	InitContextData()
	BindParam(obj interface{}) error
	Response(msgCode int, data interface{}, err error)
}

type Context struct {
	*gin.Context
}

func newContext(c *gin.Context) *Context {
	if c == nil {
		return nil
	}
	newCtx := &Context{Context: c}
	newCtx.InitContextData()
	return newCtx
}

func (c *Context) InitContextData() {
	// 设置header数据到上下文
	requestID := c.GetHeader(HEADER_REQUEST_ID)
	if requestID == "" {
		requestID = tools.GetRandomString(3)
	}
	c.Set(CONTEXT_VALUE_REQUEST_ID, requestID)
	c.Set(CONTEXT_VALUE_LOGIN_TOKEN, c.GetHeader(HEADER_LOGIN_TOKEN))
	c.Set(CONTEXT_VALUE_OV_TOKEN, c.GetHeader(HEADER_OV_TOKEN))
	c.Set(CONTEXT_VALUE_DEVICE_TYPE, c.GetHeader(HEADER_DEVICE_TYPE))
	c.Set(CONTEXT_VALUE_DEVICE_ID, c.GetHeader(HEADER_DEVICE_ID))
	c.Set(CONTEXT_VALUE_FAMILY_ID, c.GetHeader(HEADER_FAMILY_ID))
}

func (c *Context) GetRequestId() string {
	requestID := c.GetString(CONTEXT_VALUE_REQUEST_ID)
	if requestID == "" {
		requestID = tools.GetRandomString(3)
		c.Set(CONTEXT_VALUE_REQUEST_ID, requestID)
	}
	return c.GetString(CONTEXT_VALUE_REQUEST_ID)
}

func (c *Context) GetUserID() string {
	return c.GetString(CONTEXT_VALUE_UUID)
}

func (c *Context) GetFamilyID() string {
	return c.GetString(CONTEXT_VALUE_FAMILY_ID)
}

func (c *Context) GetDeviceID() string {
	return c.GetString(CONTEXT_VALUE_DEVICE_ID)
}

func (c *Context) IsIOS() bool {
	deviceType := c.GetDeviceType()
	return deviceType == DEVICE_TYPE_IOS
}

func (c *Context) IsAndroid() bool {
	deviceType := c.GetDeviceType()
	return deviceType == DEVICE_TYPE_ANDROID
}

func (c *Context) IsRobot() bool {
	deviceType := c.GetDeviceType()
	return deviceType == DEVICE_TYPE_ROBOT
}

func (c *Context) IsAPP() bool {
	deviceType := c.GetDeviceType()
	return (deviceType == DEVICE_TYPE_IOS || deviceType == DEVICE_TYPE_ANDROID)
}

func (c *Context) GetDeviceType() string {
	return c.GetString(CONTEXT_VALUE_DEVICE_TYPE)
}

func (c *Context) GetRequestError() error {
	if err, has := c.Get(CONTEXT_VALUE_REQUEST_ERR); has && err != nil {
		if requestErr, ok := err.(error); ok {
			return requestErr
		}
	}
	return nil
}

func (c *Context) BindParam(obj interface{}) error {
	return c.ShouldBind(obj)
}

func (c *Context) Response(msgCode int, data interface{}, err error) (msgData *message.ResponseData) {
	rid := c.GetRequestId()
	// 客户端错误400 服务端错误500  错误码规范更新
	httpStatus := message.HttpStatus(msgCode)
	if err != nil {
		if httpStatus == http.StatusInternalServerError {
			// 服务端错误500 触发飞书提醒
			log.Errorf("Request-ID:[%s] %+v", rid, err)
		} else {
			// 客户端错误400
			log.Warnf("Request-ID:[%s] %+v", rid, err)
		}
	}
	msgData = message.Response(msgCode, data, rid)
	c.JSON(http.StatusOK, msgData)
	return
}

func (c *Context) LogFormat(format string, a ...interface{}) string {
	_, file, line, _ := runtime.Caller(1)
	format = fmt.Sprintf("\r\n%s:%d\r\nRequest-ID:[%s] Log:[%s]", file, line, c.GetRequestId(), format)
	return fmt.Sprintf(format, a...)
}

func (c *Context) SetRequestError(err error) {
	c.Set(CONTEXT_VALUE_REQUEST_ERR, err)
}

func CreateTestContext(request *http.Request) *Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = request
	return newContext(c)
}

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := newContext(c)
		if ctx == nil {
			c.Next()
			return
		}
		h(ctx)
	}
}
