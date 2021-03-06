package redis

import (
	"filestore-server/config"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	pool *redis.Pool
	redisHost = config.RedisHost
	redisPass = config.ReidsPass
)

func newRedisPool() *redis.Pool {
	return &redis.Pool {
		MaxIdle: 	50,
		MaxActive:	30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error){
			// 1 打开连接
			conn, err := redis.Dial("tcp", redisHost)

			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			//2 访问认证
			if _, err = conn.Do("AUTH", redisPass); err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := conn.Do("PING")

			return err
		},
	}
}

func init(){
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
