package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Response struct {
	Success bool
	Message string
	Data    json.RawMessage
}

const port = ":80"

var (
	foobar       string = `{"Success": false,"Message": "Internal server error :(","Data": {"foo": "bar"}}`
	success      bool   = false
	responseCode int    = 500
	message      string
	apiToken     string = "zAZ7EtwfqYxJt8eKBRf9xfs8SQk3F4Hv22Wt29k6nchMDpeknGFhkMQeDhxBDEWS45E3dhkQNKTXqq97qCJeCZzEt3kkBfEPAC5X"
)

func main() {

	// Release the routher!!!
	r := httprouter.New()

	// API Root - Can also be used to ping the API for status check & info
	r.GET("/api/v1", index)
	r.POST("/api/v1", index)

	// API Endpoints (EP)
	r.POST("/api/v1/_reindex", AuthCheck(reindex))

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

		} else {

			// Bad token, respond with unauthorized
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

		} // Token check

	} // Nameless function
}

/* ----------------------------------------------------------------------------
Function used to reindex a single item.

@TODO - Unit test!!!!!
-----------------------------------------------------------------------------*/
func reindex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}
