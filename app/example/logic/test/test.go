package test

import (
	"demo-server/app/internal/context"
	"strconv"

	"github.com/pkg/errors"

	"demo-server/lib/log"
)

//logic层定义请求参数和返回
type TestParam struct {
	Arg1 string `json:"arg_1" binding:"required"`
	Arg2 string
}

type TestData struct {
	Arg1 string
	Arg2 string
}

func (p *TestParam) Test(c context.IContextAPI) (*TestData, error) {
	// 约定调用方返回error时附加上堆栈信息
	if _, err := strconv.Atoi(p.Arg1); err != nil {
		return nil, errors.WithStack(err)
		//return nil, errors.Wrap(err, "redis error")
		//return nil, errors.Errorf("Test strconv.Atoi(buf) error [%s]", "redis error")
		//return nil, errors.Errorf("Test strconv.Atoi(buf) error [%s]", err.Error())
	}
	// 调试信息可以使用公共方法，方便输出代码行号和请求ID
	// 示例
	//2020-06-08T17:18:55.707+0800    DEBUG   test/test.go:31
	///Users/milesc/go/src/demo-server/app/robot/logic/test/test.go:31
	//Request-ID:[1591607935dxo] Log:[hello word]
	log.Debug(c.LogFormat("hello %s", "word"))
	return &TestData{
		Arg1: p.Arg1,
		Arg2: p.Arg2,
	}, nil
}
