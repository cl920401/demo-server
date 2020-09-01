package redis

import (
	"errors"
	"demo-server/utils/tools"

	"demo-server/lib/log"
	_ "demo-server/lib/redis"
)

//基于Redis实现分布式锁
const lockPrefix = "redis_lock_"
const lockExpress = TTLSec * 60

var ErrUnLockInvalid = errors.New("lock id invalid")
var ErrHasLock = errors.New("listener locked")
var ErrConnectInvalid = errors.New("connect invalid")

// Lock : 加锁 expiration是超时时间，防止死锁，必须设置，否则走默认值
func (v *Connect) Lock(key string, expiration int) (string, error) {
	if expiration <= 0 {
		expiration = lockExpress
	}
	lockID := tools.GenerateGUID()
	ret, err := v.conn.Do("SET", KeyPrefix+lockPrefix+key, lockID, "NX", "EX", expiration)
	if ret == nil {
		return "", ErrHasLock
	}
	if err != nil {
		return "", err
	}
	return lockID, nil
}

// Unlock : 解锁
func (v *Connect) Unlock(lockKey, lockID string) error {
	lockValue, err := v.Get(lockPrefix + lockKey)
	if err != nil {
		return err
	}
	if tools.Bytes2str(lockValue) != lockID {
		log.Info(tools.Bytes2str(lockValue))
		log.Info(lockID)
		return ErrUnLockInvalid
	}

	//log.Info("Unlock Del:", ret)
	if _, err := v.Del(lockPrefix + lockKey); err != nil {
		return err
	}
	return nil
}
