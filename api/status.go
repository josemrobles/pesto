package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

type StatusResponseData struct {
	BatchID   string
	BatchSize int64
	Completed int
	Errors    int
}

func status(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var err error = nil
	var data json.RawMessage
	batchID := p.ByName("batch_id")
	redis, err := redisConn(os.Getenv("REDIS_CONNECTION"))

	// Verify that redis connection was made
	if err != nil {

		success = false
		responseCode = 500
		message = "Internal Error :("
		log.Printf("ERR: Could not connect to Redis %q", err)

	} else {

		check, err := redis.Do("SISMEMBER", "data:jobs", batchID)

		// Make sure there was no error in checking for the job
		if err != nil {

			success = false
			responseCode = 500
			message = "Internal Error :("
			log.Printf("ERR: Could not check Redis for job - %q", err)

		} else {

			if check == int64(1) {

				// GEt the number of jobs in the batch
				nj, err := redis.Do("HLEN", "stats:job:"+batchID)

				// Check if job pop inquiry failed
				if err != nil {

					success = false
					responseCode = 500
					message = "Internal Error :("
					log.Printf("ERR: Could not check job count for batch - %q", err)

				} else {

					numJobs := nj.(int64)

					responseData := &StatusResponseData{
						BatchID:   batchID,
						BatchSize: numJobs,
						Completed: 9999,
						Errors:    9999,
					}

					doof, _ := redis.Do("HGETALL", "stats:job:"+batchID)
					log.Println(doof)

					// JSONify the response data
					data, err = JSONify2(responseData)

					if err != nil {

						success = false
						responseCode = 500
						message = "Internal Error :("
						log.Printf("ERR: Could not jsonify response - %q", err)

					} else {

						// batch does exist
						success = true
						responseCode = 200
						message = "Batch found" // @TODO - Better message i.e. in progress, errors, completed, etc

					} // JSONify response

				} // Job count check

			} else {

				// batch does not exist
				success = false
				responseCode = 204
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
