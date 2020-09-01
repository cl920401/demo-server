package server

import (
	"bytes"
	"demo-server/lib/apm"
	"demo-server/lib/catch"
	"demo-server/lib/log"
	"demo-server/lib/util"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var startTime = time.Now()
var readyToServe = false
var defaultSkipPaths = []string{"/", "/sys/health", "/favicon.ico"}

//API 定义一个应用
type API struct {
	Engine *gin.Engine
}

//RegisterAPIRouter 注册路由规则
func (v *API) RegisterAPIRouter(fun func(router *gin.Engine)) {
	fun(v.Engine)
}

// NewWith new app using the specified print log skippaths
func NewWith(skipPaths []string, middleware ...gin.HandlerFunc) *API {
	if length := len(skipPaths); length == 0 {
		skipPaths = defaultSkipPaths
	}

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skipPaths,
	}), Debug(skipPaths))
	return NewApp(router, middleware...)
}

//New new api
func New(middleware ...gin.HandlerFunc) *API {

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: defaultSkipPaths,
	}), Debug(defaultSkipPaths))
	return NewApp(router, middleware...)
}

func NewApp(router *gin.Engine, middleware ...gin.HandlerFunc) *API {

	router.Use(record(), catch.Gin())
	router.Use(middleware...)
	router.Any("/sys/health", health)
	router.Any("/sys/ready", ready)
	router.Any("/", health)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
	})
	return &API{
		Engine: router,
	}
}

func record() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

//Debug debug mode
func Debug(skipPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if util.InSlice(c.Request.URL.Path, skipPaths) {
			c.Next()
			return
		}

		// TODO 2020/5/28 17:44 kangyunjie 上传文件时，未打印日志
		var requestId = c.GetHeader("Request_id")
		if requestId == "" {
			requestId, _ = uuid.GenerateUUID()
		}
		if !strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data;") {
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			}
			log.Debugf("[%s] api request, clientIp: %s, method: %s, url: %s, headers: %s, params: %s, body: %s", requestId, c.ClientIP(),
				c.Request.Method, c.Request.URL.Path, c.Request.Header, c.Request.URL.RawQuery, string(bodyBytes))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if fullPath := c.FullPath(); fullPath != "" {
			s := time.Now()
			counter := apm.Counter(fullPath, "current") //当前请求量
			//总请求计数
			if meter := apm.Meter(fullPath, "total"); meter != nil {
				meter.Mark(1)
			}
			if counter != nil {
				counter.Inc(1)
			}
			blw := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
			c.Next()
			log.Debugf("[%s] api response: %v, cost: %v", requestId, blw.body.String(), time.Since(s).Milliseconds())
			if counter != nil {
				counter.Dec(1)
			}
			if status := apm.Meter(fullPath, strconv.Itoa(c.Writer.Status())); status != nil {
				status.Mark(1)
			}

			apm.Histograms(fullPath, "execTime").Update(time.Since(s).Milliseconds())
		}
	}
}

//health 健康检查
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"pid":        os.Getpid(),
		"start_time": startTime,
		"code":       http.StatusOK,
	})
}

func ready(c *gin.Context) {
	if readyToServe {
		c.JSON(http.StatusOK, gin.H{
			"pid":        os.Getpid(),
			"start_time": startTime,
			"code":       http.StatusOK,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"pid":        os.Getpid(),
			"start_time": startTime,
			"code":       http.StatusBadRequest,
		})
	}
}

func SetReadyToServe(isReady bool) {
	readyToServe = isReady
}
