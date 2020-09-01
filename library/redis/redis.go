package redis

import (
	"demo-server/lib/log"

	"github.com/gomodule/redigo/redis"
)

// 定义基础配置
const (
	TTLSec    = 1
	KeyPrefix = "demo_server_"
)

// Connect : redis连接
type Connect struct {
	conn redis.Conn
}

// Conn : 获取redis可用连接
func Conn(key string) *Connect {
	if !isInit {
		rw.Lock()
		defer rw.Unlock()
		InitRedis()
	}
	if r, ok := conn[key]; ok && r != nil {
		return &Connect{conn: r.Get()}
	}
	return nil
}

// Default : 获取默认实例
func Default() *Connect {
	return Conn("default")
}

// Redis基本操作封装

// SetExString : 设置键值对
func (v *Connect) SetExString(key, value string, seconds interface{}) (reply interface{}, err error) {
	return v.conn.Do("SETEX", KeyPrefix+key, seconds, value)
}

// HSetExString : 设置多字段键值对
func (v *Connect) HSetExString(key, field, value string) (reply interface{}, err error) {
	return v.conn.Do("hset", KeyPrefix+key, field, value)
}

// HMSetExString :
func (v *Connect) HMSetExString(key string, data map[string]string) (reply interface{}, err error) {
	return v.conn.Do("hmset", redis.Args{}.Add(KeyPrefix+key).AddFlat(data)...)
}

func (v *Connect) Exist(key string) (bool, error) {
	return redis.Bool(v.conn.Do("EXISTS", KeyPrefix+key))
}

// Get : 读取键值对
func (v *Connect) Get(key string) ([]byte, error) {
	return redis.Bytes(v.conn.Do("GET", KeyPrefix+key))
}

// GetString : 读取字符串格式键值对
func (v *Connect) GetKeys(key string) ([]string, error) {
	val, err := v.conn.Do("KEYS", KeyPrefix+key)
	return redis.Strings(val, err)
}

// DelWithoutPrefix : 删除键值对(忽略Prefix)
func (v *Connect) DelWithoutPrefix(key string) (interface{}, error) {
	return v.conn.Do("DEL", key)
}

// Del : 删除键值对
func (v *Connect) Del(key string) (interface{}, error) {
	return v.conn.Do("DEL", KeyPrefix+key)
}

// Get : 读取多字段键值对
func (v *Connect) HGet(key, field string) (interface{}, error) {
	return v.conn.Do("HGET", KeyPrefix+key, field)
}

// HGetString : 读取字符串格式多字段键值对
func (v *Connect) HGetString(key, field string) (string, error) {
	val, err := v.HGet(key, field)
	return redis.String(val, err)
}

// HGetStringWithoutPrefix : 不处理key前缀
func (v *Connect) HGetStringWithoutPrefix(key, field string) (string, error) {
	val, err := v.conn.Do("HGET", key, field)
	return redis.String(val, err)
}

// HMGetString : 读取字符串格式多字段键值对
func (v *Connect) HGetAll(key string) (map[string]string, error) {
	val, err := v.conn.Do("HGETALL", KeyPrefix+key)
	return redis.StringMap(val, err)
}

func (v *Connect) Expire(key string, seconds int64) error {
	_, err := v.conn.Do("EXPIRE", key, seconds)
	return err
}

// LPeekList : 查看list第一个值
func (v *Connect) LPeekList(key string) ([]byte, error) {
	return redis.Bytes(v.conn.Do("LINDEX", KeyPrefix+key, 0))
}

// RPushList : list数据右入列
func (v *Connect) RPushList(key string, value []byte) (interface{}, error) {
	return v.conn.Do("RPUSH", KeyPrefix+key, value)
}

// LPopList : list数据左推出
func (v *Connect) LPopList(key string) ([]byte, error) {
	return redis.Bytes(v.conn.Do("LPOP", KeyPrefix+key))
}

// Close : 关闭连接
func (v *Connect) Close() {
	if err := v.conn.Close(); err != nil {
		log.Error("Connect Close error:", err)
	}
}

//func GetSameKeyData() (interface{}, error) {
//	conn := GetRedis().Get()
//	defer func() {
//		err := conn.Close()
//		if err != nil {
//			log.Error(err)
//		}
//	}()
//	return redis.Values(conn.Do("KEYS", KeyPrefix))
//}
