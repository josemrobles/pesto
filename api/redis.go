package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

/* ----------------------------------------------------------------------------
Function used to connect to the redis server...

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func connectToRedis(a string) redis.Conn {

	// Connect to redis
	c, err := redis.Dial("tcp", a)

	// Check if the connection was successful
	if err != nil {

		// Foobar, could not connect to the redis server
		log.Printf("ERR: Could not connect to redis - %q", err)
	}
	return c
}
