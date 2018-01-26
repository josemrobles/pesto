package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
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
		bID,bSize,err := processBatch(b)

		// Check for errors 
		if err != nil {

			// Foobar no wascally wabbits!!
			success = false
			responseCode = 500
			message = "Internal Error :("
			log.Printf("ERR: Could not process batch - %q", err)

		} else {


			responseData := &ResponseData{
				BatchID: bID,
				BatchCount: bSize,
			}

			// JSONify the response data
			data,err = JSONify(responseData)

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

			}

		}

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
