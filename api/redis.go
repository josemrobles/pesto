package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"log"
)


func initRedis(a string, connectTimeout, readTimeout, writeTimeout time.Duration) (redis.Conn, error) {
	return redis.Dial("tcp", a,
		redis.DialConnectTimeout(connectTimeout),
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout))
}


/* ----------------------------------------------------------------------------
Function used to connect to the redis server...

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func initRedisPool(a string,idleTimeout time.Duration) (redis.Pool, error) {

	// INitialize the Redis pool
	p := redis.Pool{
		MaxIdle:     10,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {

			// Connect to Redis
			c, err := initRedis(a,10*time.Second ,10*time.Second ,10*time.Second)

			// Check if we were able to connect
			if err != nil {

				log.Printf("ERR: Unable to connect to redis server - %q",err)
				return nil, err

			} else {
				log.Println("OOF: CARALHO!!!!!!!!!!!!!!!!!!!!!!!!!")
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return p, nil
}
