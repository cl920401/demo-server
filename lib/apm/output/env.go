package output

import (
	"demo-server/lib/config"
	"os"
	"strings"
)

func ShortHostname() string {
	host, _ := os.Hostname()
	if index := strings.Index(host, "."); index > 0 {
		return host[:index]
	}
	return host
}

func ServiceName() string {
	return config.Get("service.name").String("")
}
