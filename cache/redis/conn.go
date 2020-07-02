package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool *redis.Pool
	redisHost = "127.0.0.1"
)

func newRedisPool() *redis.Pool  {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				println("redis链接失败,err=" + err.Error())
				return nil, err
			}
			return c, nil
		},
		DialContext:     nil,
		//定时检查redis连接是否可用
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         50,
		MaxActive:       30,
		IdleTimeout:     300 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}

func init()  {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}