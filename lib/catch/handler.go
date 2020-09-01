package catch

import (
	"context"
	"demo-server/lib/apm"
	"demo-server/lib/config"
	"demo-server/lib/feishu"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"runtime/debug"
)

func GRPC(_ context.Context, p interface{}) (err error) {
	enable := config.Get("feishu.enable").Bool(false)
	if enable {
		send(fmt.Sprintf("%s", p), string(debug.Stack()))
	}
	return status.Errorf(codes.Internal, "%s", p)
}

func Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				send(fmt.Sprintf("%s", p), string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}()
		c.Next()
	}
}

func send(title, text string) {
	feishu.Send(title, text)
	if meter := apm.Meter("panic", "total"); meter != nil {
		meter.Mark(1)
	}
}
