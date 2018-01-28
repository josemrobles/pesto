package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"net/http"
)

type Status struct {
	BatchID   string
	BatchSize int64
	Completed int
	Errors    int
}

/* ----------------------------------------------------------------------------
Function used to obtain the sum of a specified status for a specified batch.
Statuses include 0 = in progress, 1 =  done, 2 = error.

@TODO - Unit test!!!!!
@TODO - better error handling
-----------------------------------------------------------------------------*/
func getCountByStatus(b string,s int)  (r int){

	// Grab redis connection from redis pool
	c := redisPool.Get()
	defer c.Close()

	defer func() {
		c.Close()
	}()

	// Get all jobs for the specified batch
	m, err := redis.StringMap(c.Do("hgetall", "stats:job:"+b))

	// Check to make sure we are able to get the data
	if err != nil {

		log.Printf("ERR: Could not get failed job count - %q", err)

	}  else {

		// Iterate through the has and count the number of hits.
		for _, v := range m {

			// Convert the value to int for comparison
			vi,_ := strconv.Atoi(v)

			if vi == s {
				r++
			} // If

		} // Loop

	} // Check

    return
}

func status(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var err error = nil
	var data json.RawMessage
	batchID := p.ByName("batch_id")
	c := redisPool.Get()
	defer c.Close()

	defer func() {
		c.Close()
	}()

	check, err := c.Do("SISMEMBER", "data:jobs", batchID)

	// Make sure there was no error in checking for the job
	if err != nil {

		success = false
		responseCode = 500
		message = "Internal Error :("
		log.Printf("ERR: Could not check Redis for job - %q", err)

	} else {

		if check == int64(1) {

			// GEt the number of jobs in the batch
			nj, err := c.Do("HLEN", "stats:job:"+batchID)

			// Check if job pop inquiry failed
			if err != nil {

				success = false
				responseCode = 500
				message = "Internal Error :("
				log.Printf("ERR: Could not check job count for batch - %q", err)

			} else {

				numJobs := nj.(int64)

				status := &Status{
					BatchID:   batchID,
					BatchSize: numJobs,
					Completed: getCountByStatus(batchID,1),
					Errors:    getCountByStatus(batchID,2),
				}

				// JSONify the response data
				data, err = JSONify2(status)

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
