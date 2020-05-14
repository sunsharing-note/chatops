package myredis

import "github.com/gomodule/redigo/redis"

// 初始化redis
var MyPool *redis.Pool

func RedisPool()*redis.Pool{
	return &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp","122.51.79.172:6379")
		},
		TestOnBorrow:    nil,
		MaxIdle:         20,
		MaxActive:       0,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}
