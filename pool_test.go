package redix

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestPool(t *testing.T) {
	src := newPool("localhost:6378", "", 0)
	dest := newPool("localhost:6379", "", 0)
	pool := NewRedisPoolx(
		src,
		dest,
		func() string {
			return "dest"
		},
	)

	conn := pool.Get()
	defer conn.Close()

// 	conn.Do("set", "kv-k", "hello")

	rs, e := redis.String(conn.Do("get", "kv-k"))
	if e != nil {
		panic(e)
	}

	fmt.Println(rs)
}
