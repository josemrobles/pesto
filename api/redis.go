package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"log"
)

/* ----------------------------------------------------------------------------
Function used to connect to the redis server...

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func initRedisPool(a string) *redis.Pool {

	// INitialize the Redis pool
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 120*time.Second,
		Dial: func() (redis.Conn, error) {

			// Connect to Redis
			c, err := redis.Dial("tcp", a)

			// Check if we were able to connect
			if err != nil {

				log.Printf("ERR: Unable to connect to redis server - %q",err)
				return nil, err

			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
