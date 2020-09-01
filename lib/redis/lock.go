package redis

import (
	"context"
	"demo-server/lib/log"
	rds "github.com/gomodule/redigo/redis"
	"time"
)

const expire = 3

/*
Lock
key: redis lock key
try: 重试次数
*/
func Lock(key string, try int8) (unLock func(), ok bool) {

	for i := 0; i <= int(try); i++ {
		conn := Default()
		if ok, err := rds.String(conn.Do("SET", key, "lock", "ex", expire, "nx")); err == nil && ok == "OK" {
			conn.Close()
			//设置锁成功 启动后台协程维护ttl
			ctx, cancel := context.WithCancel(context.Background())
			go start(ctx, key, expire)

			return func() {
				cancel() //去掉后台协程
				conn := Default()
				//删除redis key
				if err := conn.Del(key); err != nil {
					log.Error(err)
				}
			}, true
		}
		conn.Close()
		if try > 0 {
			time.Sleep(1 * time.Second)
		}
	}
	return nil, false
}

func start(ctx context.Context, key string, ttl int64) {
	tick := time.Tick(time.Duration(ttl) / 2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			conn := Default()
			if err := conn.Expire(key, ttl); err != nil {
				log.Error(err)
			}
		}
	}
}
