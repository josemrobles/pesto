package main

import (
	"encoding/json"
	"fmt"
	"github.com/josemrobles/conejo"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Success bool
	Message string
	Data    json.RawMessage
}

type ResponseData struct {
	BatchID string
	BatchCount int
}

const port = ":80"

var (
	rmq                 = conejo.Connect(os.Getenv("RABBITMQ_CONNECTION"))
	queue               = conejo.Queue{Name: os.Getenv("RABBITMQ_QUEUE"), Durable: false, Delete: false, Exclusive: false, NoWait: false}
	exchange            = conejo.Exchange{Name: os.Getenv("RABBITMQ_EXCHANGE"), Type: "topic", Durable: true, AutoDeleted: false, Internal: false, NoWait: false}
	foobar       string = `{"Success": false,"Message": "Internal server error :(","Data": {"foo": "bar"}}`
	success      bool   = false
	responseCode int    = 500
	message      string
	data         json.RawMessage
	apiToken     string = string(os.Getenv("API_TOKEN"))
)

func main() {

	// Release the routher!!!
	r := httprouter.New()

	// API Root - Can also be used to ping the API for status check & info
	r.GET("/api/v1", index)
	r.POST("/api/v1", index)

	// API Endpoints (EP)
	r.POST("/api/v1/_reindex", AuthCheck(reindex))

	// Caralho, it no chooch!
	log.Fatal(http.ListenAndServe(port, r))
}

/* ----------------------------------------------------------------------------
API Index function used as a general health check endpoint i.e. ping | pong.
Should be used for app monitoring.

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Prep the API response
	ping := &Response{
		Success: true,
		Message: "pong",
		Data:    json.RawMessage(`{}`),
	}

	// Marshal the response in preparation for output
	pong, err := json.Marshal(ping)

	// Check if there was an error in the Marshal for pong
	if err != nil {

		// Fubar, could not marshal the response
		log.Println("ERR: Could not Marshal ping - [ index ]")
		fmt.Fprint(w, "ERROR")

	} else {

		// All is well, reply to the png
		fmt.Fprint(w, string(pong))

	} // Marshall check
}

/* ----------------------------------------------------------------------------
API middleware used to validate the incoming request and add anny additional
logging or metrics for future analysis.

@TODO - Unit test!!!!!
@TODO - Proper auth token check
-----------------------------------------------------------------------------*/
func AuthCheck(h httprouter.Handle) httprouter.Handle {

	// A function has no name...
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// Check the provided token against the token var.
		if string(r.Header["Token"][0]) == apiToken {

			// Valid token, move along
			h(w, r, ps)
			log.Println("OOF: Request accepted.")

		} else {

			// Bad token, respond with unauthorized
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Println("OOF: Request submitted with invalid token.")

		} // Token check

	} // Nameless function
}

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
				message = "Payload Accepted, check status of the request ["+bID +"]"

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
