package router

import (
	"demo-server/app/internal/context"
)

//Info 路由信息
type Info struct {
	Method            string
	Handler           context.HandlerFunc
	Name              string
	PermissionCheck   bool
	PermissionCheckV2 bool
}
