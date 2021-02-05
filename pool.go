package redix

import (
	"github.com/garyburd/redigo/redis"
)


type RedisPoolI interface{
	Get() redis.Conn
}

type RedisPoolx struct {
	readWho  string      // src 读源， dest 读目标
	srcPool  *redis.Pool // 数据源切换的A方
	destPool *redis.Pool // 数据源切换的B方
}

func NewRedisPoolx(srcPool *redis.Pool, destPool *redis.Pool, getReadWho func() string) *RedisPoolx {
	var readWho string
	if getReadWho == nil {
		readWho = "src"
	} else {
		readWho = getReadWho()

		if readWho == "" {
			readWho = "src"
		}
	}

	return &RedisPoolx{
		readWho:  readWho,
		srcPool:  srcPool,
		destPool: destPool,
	}
}

func (rp *RedisPoolx) Get() redis.Conn {
	srcConn := rp.srcPool.Get()
	destConn := rp.destPool.Get()

	return NewConnx(srcConn, destConn, rp.readWho)
}
