package main

import (
	"encoding/json"
	"fmt"
	"github.com/josemrobles/conejo"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/* ----------------------------------------------------------------------------
Function used to reindex a single item.

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func reindex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Decode the payload
	b, err := ioutil.ReadAll(r.Body)

	// Check if we were able to read the payload
	if err != nil {

		// Meh - Could not decode the payload...fml
		success = false
		responseCode = 500
		message = "Internal Error :("
		log.Printf("ERR: Could not read POST data - %q", err)

	} else {

		// Process the batch / request
		batchID, batchSize, err := processBatch(b)

		// Check for errors
		if err != nil {

			// Foobar no wascally wabbits!!
			success = false
			responseCode = 500
			message = "Internal Error :("
			log.Printf("ERR: Could not process batch - %q", err)

		} else {

			responseData := &ResponseData{
				BatchID:    batchID,
				BatchCount: batchSize,
			}

			// JSONify the response data
			data, err = JSONify(responseData)

			if err != nil {

				// Foobar no wascally wabbits!!
				success = false
				responseCode = 500
				message = "Internal Error :("
				log.Printf("ERR: Could not process batch - %q", err)

			} else {

				success = true
				responseCode = 202 // Accepted
				message = "Request accepted"

			} // JSONification

		} // processBatch()

	} // Read payload

	// By this point we should have some sort of response
	resp := &Response{Success: success, Message: message, Data: data}

	// SET content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the response
	response, err := json.Marshal(resp)

	// Check to see if there was an error whilst marshalling the response
	if err != nil {

		// FML
		log.Printf("ERR: Could not marshal response - %q", err)
		w.WriteHeader(500)
		fmt.Fprint(w, foobar)

	} else {

		// Respond
		w.WriteHeader(responseCode)
		fmt.Fprint(w, string(response))
	}
}

/* ----------------------------------------------------------------------------
Used to generate the unique batch ID that can be used to lookup the status of
the batch / job via the /status API endpoint.

@TODO - Unit test!!!!!
@TODO - Make the hash length configurable? via func param?
-----------------------------------------------------------------------------*/
func getBatchID() string {

	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("12345678abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 50)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

/* ----------------------------------------------------------------------------
Function used to process the incoming payload. This func will iterate through
the payload and break it up into individual jobs which can be tracked individually
to obtain a granular batch status report.

@TODO - Unit test!!!!!
@TODO - Actually iterate through the batch
@TODO - Add the individual payload to the redis hash
@TODO - Should I create the batch before or after I publish the payload(s)?
-----------------------------------------------------------------------------*/
func processBatch(b []byte) (string, int, error) {

	var err error = nil
	c := redisPool.Get()
	defer c.Close()

	// Get new batch ID
	batchID := getBatchID()

	// Get total number of jobs in batch
	numJobs := 10

	// Add new batch to redis
	c.Do("SADD", "data:jobs", batchID)

	// Iterate through the payload and send each message
	// @TODO - Actually iterate through the payload, currently a simulation
	for i := 0; i < numJobs; i++ {

		// Convert item to string
		item := strconv.Itoa(i + 1)

		// Set the status for the current job 0 = processing 1 = done 2 = error
		c.Do("HSET", "stats:job:"+batchID, "job:"+item+":status", 0)

		// Publish the message
		err = conejo.Publish(rmq, queue, exchange, string([]byte(b)))

		// Check to make sure the there were no errors in publishing
		if err != nil {

			log.Printf("ERR: Could not publish message %v - %q", i, err)

		} // Publish message

	} // Iterate

	return batchID, numJobs, err

}
