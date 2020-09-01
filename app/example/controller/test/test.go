package test

import (
	logic "demo-server/app/example/logic/test"
	"demo-server/app/internal/context"
	"demo-server/app/internal/message"

	"github.com/pkg/errors"
)

// Test : test
func Test(c *context.Context) *message.ResponseData {
	var p logic.TestParam
	if err := c.BindParam(&p); err != nil {
		return c.Response(message.CodeParamErr, nil, errors.Errorf("Test c.BindParam(&p) error [%s]", err.Error()))
	}
	data, err := p.Test(c)
	if err != nil {
		return c.Response(message.CodeUnknowErr, nil, errors.WithMessage(err, "p.Test(c.Context())"))
	}

	return c.Response(message.CodeOK, data, nil)
	//这里不需要return，代码检查工具会警告
}
