package main

import (
	"github.com/garyburd/redigo/redis"
)

/* ----------------------------------------------------------------------------
Function used to connect to the redis server...

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func redisConn(a string) (redis.Conn, error) {

	// Connect to redis
	c, err := redis.Dial("tcp", a)

	return c, err
}

func initPool(ip string, maxIdleConnection int, idleTimeout, connectTimeout, readTimeout, writeTimeout time.Duration) (redis.Pool, error) {

	pool := redis.Pool{
		MaxIdle:     maxIdleConnection,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := initRedis("tcp", ip, connectTimeout, readTimeout, writeTimeout)
			if err != nil {
				log.Errorf("connect error", conn, err)
				return nil, err
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return pool, nil
}

func initRedis(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (redis.Conn, error) {
	return redis.Dial(network, address,
		redis.DialConnectTimeout(connectTimeout),
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout))
}
