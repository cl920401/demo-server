// +build windows

package server

import (
	"demo-server/lib/log"
	"net/http"
)

//Run 运行
func (v *API) Run(addr string) {
	log.Info("server start and listening on:", addr)
	srv := &http.Server{Addr: addr, Handler: v.Engine}
	err := srv.ListenAndServe()

	if err != nil {
		log.Error(err)
	}
}
