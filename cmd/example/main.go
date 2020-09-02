package main

import (
	//自动设置GOMAXPROCS，用于处理运行在容器中CPU quota
	_ "go.uber.org/automaxprocs"

	"demo-server/app/example/router"
	"demo-server/lib/mysql"
	"demo-server/lib/redis"
	"demo-server/lib/server"
	"demo-server/middleware"

	"github.com/urfave/cli"

	"fmt"
	"os"
	"sort"
	"time"

	conf "demo-server/lib/config"
	"demo-server/lib/log"
)

const appName = "example"
const appVersion = "2.0.0804"

func main() {
	app := cli.NewApp()
	app.Action = run
	app.Name = appName
	version = appVersion
	app.Compiled = time.Now()
	app.Version = fmt.Sprintf("%s\n branch: %s\n commit: %s\n compileAt:%s", version, branch, commit, app.Compiled)

	app.Usage = fmt.Sprintf("The %s command line interface.", appName)
	app.Copyright = "Copyright 2019-2020 The Authors."

	app.Flags = append(app.Flags, &ConfigFlag, &RPCPortFlag)
	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
	fmt.Print("init project")
}

func run(ctx *cli.Context) error {

	err := conf.LoadFile(ctx.String("config"))
	if err != nil {
		log.Fatalf("config file non existent: %s", ctx.String("config"))
	}
	log.Debug(conf.ToString())

	redis.InitRedis()
	mysql.InitDB()
	//启动服务
	app := server.New(middleware.Cors(), middleware.Hook())
	app.RegisterAPIRouter(router.API)
	app.Run(conf.Get("port").String("8081"))
	return nil
}
