package redis

import (
	"demo-server/lib/log"
	"demo-server/lib/config"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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
type Connect struct {
	redis.Conn
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

//Conn  获取redis可用连接
func Conn(key string) *Connect {
	if !isInit {
		rw.Lock()
		defer rw.Unlock()
		InitRedis()
	}
	if r, ok := conn[key]; ok && r != nil {
		return &Connect{
			Conn: r.Get(),
		}
	}
	return nil
}

//Default 获取默认实例
func Default() *Connect {
	return Conn("default")
}

//InitRedis init redis config
func InitRedis() {
	if isInit {
		return
	}
	dbConf := config.Get("redis")
	var data map[string]*Conf
	if err := dbConf.Scan(&data); err != nil {
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
func (conn Connect) Get(key string) ([]byte, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return redis.Bytes(conn.Do("GET", key))
}

func (conn Connect) GetBool(key string) (bool, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return redis.Bool(conn.Do("GET", key))
}

func (conn Connect) Del(key string) error {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	_, err := conn.Do("DEL", key)
	if err != nil && err != redis.ErrNil {
		log.Error("redis err:", err)
	}
	return err
}

func (conn Connect) SetBool(key string, value interface{}) (reply interface{}, err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return conn.Do("SET", key, value)
}

func (conn Connect) Set(key string, value interface{}) (reply interface{}, err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	data, _ := json.Marshal(value)
	return conn.Do("SET", key, data)
}

func (conn Connect) SetEx(key string, seconds, value interface{}) (reply interface{}, err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	data, _ := json.Marshal(value)
	return conn.Do("SETEX", key, seconds, data)
}

func (conn Connect) SetExStr(key string, seconds, value interface{}) (reply interface{}, err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	data, _ := json.Marshal(value)
	return conn.Do("SETEX", key, seconds, data)
}
func (conn Connect) SetExString(key, value string, seconds interface{}) (reply interface{}, err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return conn.Do("SETEX", key, seconds, value)
}

// Deprecated: Use Lock
func (conn Connect) Lock(key string, seconds interface{}, value string) (reply interface{}, err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return conn.Do("SET", "LOCK"+key, value, true, "ex", seconds, "nx")
}

// Deprecated: Use Lock
func (conn Connect) UnLock(key string) error {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	_, err := conn.Do("DEL", "LOCK"+key)
	if err != nil && err != redis.ErrNil {
		log.Error("redis err:", err)
	}
	return err
}

func (conn Connect) MGet(keys ...interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	args := make([]interface{}, 0, len(keys))
	args = append(args, keys...)
	return conn.Do("MGET", args...)
}

func (conn Connect) MSet(keys []interface{}, values []interface{}) (reply interface{}, err error) {
	if len(keys) != len(values) {
		return nil, fmt.Errorf("kv not match")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	args := make([]interface{}, 0, len(keys)+len(values))
	for i := 0; i < len(keys); i++ {
		args = append(args, keys[i])
		args = append(args, values[i])
	}
	return conn.Do("MSET", args...)
}

//************* LIST ****************/

func (conn Connect) RPush(key interface{}, values ...interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	args = append(args, values...)
	return conn.Do("RPUSH", args...)
}

func (conn Connect) LPush(key interface{}, values ...interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	args = append(args, values...)
	return conn.Do("LPUSH", args...)
}

func (conn Connect) LRemove(key interface{}, number int64, value interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	args := make([]interface{}, 0, 3)
	args = append(args, key)
	args = append(args, number)
	args = append(args, value)
	return conn.Do("LREM", args...)
}

func (conn Connect) LPop(key interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	return conn.Do("LPOP", key)
}

func (conn Connect) LRange(key interface{}, start, stop int64) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	// redigo parse reply
	return conn.Do("LRANGE", key, start, stop)
}

func (conn Connect) LLen(key interface{}) (int64, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return redis.Int64(conn.Do("LLEN", key))
}

//************** SET ***************

func (conn Connect) SAdd(key interface{}, members ...interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	args := make([]interface{}, 0, len(members)+1)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int64(conn.Do("SADD", args...))
}

func (conn Connect) SRemove(key interface{}, members ...interface{}) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	args := make([]interface{}, 0, len(members)+1)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int64(conn.Do("SREM", args...))
}

func (conn Connect) SIsMember(key interface{}, member interface{}) (bool, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	val, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		log.Error(err)
	}
	return val == 1, err
}

func (conn Connect) Exists(key interface{}) (bool, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	val, err := redis.Int(conn.Do("EXISTS", key))
	return val == 1, err
}

func (conn Connect) Expire(key interface{}, seconds int64) error {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, seconds)
	return err
}

// Pipeline
func (conn Connect) MExpire(keys []interface{}, seconds []int64) error {
	if len(keys) != len(seconds) {
		return fmt.Errorf("kv not match")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for i := 0; i < len(keys); i++ {
		err := conn.Send("EXPIRE", keys[i], seconds[i])
		if err != nil {
			log.Error(err)
			continue
		}
	}
	if err := conn.Flush(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (conn Connect) HSet(key string, field, value string) (interface{}, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return conn.Do("HSET", key, field, value)
}

func (conn Connect) HGet(key string, field string) (string, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return redis.String(conn.Do("HGET", key, field))
}

func (conn Connect) HDel(key string, fields ...interface{}) (int, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	if len(fields) == 0 {
		return 0, nil
	}
	args := make([]interface{}, 0, len(fields))
	args = append(args, key)
	args = append(args, fields...)
	return redis.Int(conn.Do("HDEL", args...))
}

func (conn Connect) HExists(key string, field string) (bool, error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	return redis.Bool(conn.Do("HEXISTS", key, field))
}

func (conn Connect) PubSub(channel string, handler func(data redis.Message)) error {
	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(channel); err != nil {
		return err
	}

	defer psc.Close()
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			log.Debugf("%s: message: %s", v.Channel, v.Data)
			if handler != nil {
				handler(v)
			}

		case redis.Subscription:
			log.Debugf("redis %s: %s %d", v.Channel, v.Kind, v.Count)

		case error:
			return v
		}
	}
}
