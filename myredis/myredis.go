package myredis

import (
	"github.com/gomodule/redigo/redis"
)


func RedisPool(ip string)*redis.Pool{
	return &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp",ip)
		},
		TestOnBorrow:    nil,
		MaxIdle:         20,
		MaxActive:       0,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}
