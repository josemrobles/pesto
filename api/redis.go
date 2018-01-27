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
