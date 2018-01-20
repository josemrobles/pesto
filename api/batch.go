package main

import (
	"log"
	"github.com/josemrobles/conejo"
)

func processBatch(b []byte) (string ,error) {

	var err error = nil
	var errors int = 0

	// Iterate through the payload and send each message
	// @TODO - Actually iterate through the payload, currently a simulation
	for i := 0; i < 10; i++ {

		// Publish the message
		err = conejo.Publish(rmq, queue, exchange, string([]byte(b)))

		// Check to make sure the there were no errors in publishing
		if err != nil {

			log.Printf("ERR: Could not publish message %v - %q", i,err)
			errors++

		} // Publish message
	}

	return "76589878687980897890hh9800", err

}
