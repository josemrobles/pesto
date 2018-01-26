package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"fmt"
	"encoding/json"
	"net/http"
	"os"
)

func status(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var err error = nil
	var data json.RawMessage
	batchID := p.ByName("batch_id")
	redis,err := redisConn(os.Getenv("REDIS_CONNECTION"))

	// Verify that redis connection was made
	if err != nil {

		success = false
		responseCode = 500
		message = "Internal Error :("
		log.Printf("ERR: Could not connect to Redis %q",err)

	} else {

		check, err := redis.Do("SISMEMBER", "data:jobs", batchID)

		// Make sure there was no error in checking for the job
		if err != nil {

			success = false
			responseCode = 500
			message = "Internal Error :("
			log.Printf("ERR: Could not check Redis for job - %q",err)

		} else {

			if check == int64(1) {

				// GEt the number of jobs in the batch
				numJobs,err := redis.Do("HLEN", "stats:job:"+batchID)

				// Check if job pop inquiry failed
				if err != nil {

					success = false
					responseCode = 500
					message = "Internal Error :("
					log.Printf("ERR: Could not check job count for batch - %q",err)

				} else {

					log.Println(numJobs)

				} // Job count check

			} else {

				// batch does not exist 
				success = false
				responseCode = 412
				message = "Batch not found"

			} // Batch check 

		} // job check in redis

	} // Redis connect

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
