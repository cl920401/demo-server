package redis

import (
	"fmt"
	"sync"
	"time"

	"demo-server/lib/config"
	"demo-server/lib/log"

	"github.com/gomodule/redigo/redis"
)

//Conf redis config
type Conf struct {
	Host        string `json:"host"`
	Password    string `json:"password"`
	DB          int    `json:"db"`
	Port        int16  `json:"port"`
	MaxIdle     int    `json:"max_idle"`
	MaxActive   int    `json:"max_active"`
	IdleTimeout int    `json:"idle_timeout"`
}

var isInit bool
var rw sync.RWMutex

//ToString conf to string
func (v *Conf) ToString() string {
	return fmt.Sprintf("%+v\n", v)
}

//GetAddr 格式化配置文件中地址和端口
func (v *Conf) GetAddr() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

var conn = make(map[string]*redis.Pool)

//InitRedis init redis config
func InitRedis() {
	if isInit {
		return
	}
	redisConf := config.Get("redis")
	var data map[string]*Conf
	if err := redisConf.Scan(&data); err != nil {
		log.Fatal("error parsing redis configuration file ", err)
		return
	}
	for k, v := range data {
		redisPool := newRedis(v)
		//测试是否连通
		redisConn := redisPool.Get()
		if err := redisPool.TestOnBorrow(redisConn, time.Now()); err != nil {
			log.Fatal("initRedis ERROR", k, v.ToString(), err)
		}
		if err := redisConn.Close(); err != nil {
			log.Error(err)
		}
		conn[k] = redisPool
	}
	isInit = true
}

func newRedis(conf *Conf) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.GetAddr(), redis.DialDatabase(conf.DB), redis.DialPassword(conf.Password))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, _ time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

//
//func GetRedis() *redis.Pool {
//	return pool
//}
