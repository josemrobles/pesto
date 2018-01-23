package main

import (
	"log"
	"github.com/josemrobles/conejo"
	"encoding/json"
)

func processBatch(b []byte) (string ,int,error) {

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

	return "76589878687980897890hh9800", 10,err

}

func JSONify(responseData *ResponseData) (json.RawMessage, error) {

	// Marahal the incoing response
	b, err := json.Marshal(responseData)

	// Check for an error
	if err != nil {

		// No bueno
		return nil, err

	} else {

		// Return the struct in raw json
		return json.RawMessage(string(b)), nil

	}
}
