// +build !windows

package server

import (
	"context"
	"demo-server/lib/config"
	"demo-server/lib/log"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
	_ "net/http/pprof"
	"time"
)

//Run 运行
func (v *API) Run(addr string) {
	log.Info("server start and listening on:", addr)
	go pprof()
	err := gracehttp.Serve(
		&http.Server{Addr: addr, Handler: v.Engine},
	)
	if err != nil {
		log.Error(err)
	}
}

func pprof() {
	var server *http.Server
	var isRun bool
	var port = 6020
	//检测是否需要关闭
	go func() {
		for range time.Tick(3 * time.Second) {
			if !config.Get("pprof.enable").Bool(false) && isRun && server != nil {
				_ = server.Shutdown(context.Background())
				isRun = false
				port = 6020
			}
		}
	}()

	for range time.Tick(3 * time.Second) {
		if config.Get("pprof.enable").Bool(false) && !isRun {
			addr := fmt.Sprintf("0.0.0.0:%d", port)
			log.Info("running http pprof on: ", addr)
			server = &http.Server{Addr: addr}
			isRun = true
			if err := server.ListenAndServe(); err != nil {
				isRun = false
				log.Error("running http pprof error ", err.Error())
			}
			port++
		}
	}
}
