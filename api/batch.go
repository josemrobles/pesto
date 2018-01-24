package main

import (
	"log"
	"github.com/josemrobles/conejo"
	"encoding/json"
	"os"
	"time"
	"math/rand"
)

func getBatchID() string{

	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("12345678abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 30)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func processBatch(b []byte) (string ,int,error) {

	var err error = nil
	redis,err := redisConn(os.Getenv("REDIS_CONNECTION"))

	// Get new batch ID
	bID := getBatchID()

	// Get total number of jobs in batch
	numJobs := 5

	if err != nil {

		log.Printf("ERR: Could not connect to Redis %q",err)

	} else {

		// Add new batch to redis
		redis.Do("SADD", "data:jobs", bID)
		redis.Do("HSET", "stats:job:"+bID,"total_jobs",numJobs)

		// Iterate through the payload and send each message
		// @TODO - Actually iterate through the payload, currently a simulation
		for i := 0; i < numJobs; i++ {

			// Publish the message
			err = conejo.Publish(rmq, queue, exchange, string([]byte(b)))

			// Check to make sure the there were no errors in publishing
			if err != nil {

				log.Printf("ERR: Could not publish message %v - %q", i,err)

			} else {


			} // Publish message

		} // Iterate

	} // Redis connection

	return bID, numJobs,err

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
