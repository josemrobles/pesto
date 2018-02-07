package main

import (
	"github.com/josemrobles/conejo"
	"log"
	"os"
	"runtime"
)

var (
	rmq       = conejo.Connect(os.Getenv("RABBITMQ_CONNECTION"))
	workQueue = make(chan string)
	queue     = conejo.Queue{Name: os.Getenv("RABBITMQ_QUEUE"), Durable: false, Delete: false, Exclusive: false, NoWait: false}
	exchange  = conejo.Exchange{Name: os.Getenv("RABBITMQ_EXCHANGE"), Type: "topic", Durable: true, AutoDeleted: false, Internal: false, NoWait: false}
)

/* ----------------------------------------------------------------------------
Init func launches the goroutines which do the concurrent processing of the
messages which are received via wabbitMQ. The number of goroutnes is determined
by the number of CPU's for the host / container. This allows for maximum efficiency
and makse sure the host is fully utilized. 
-----------------------------------------------------------------------------*/
func init() {

	// Lanch N goroutines based on number of cores.
	for i := 0; i < runtime.NumCPU(); i++ {

		log.Printf("OOF - launching sub-worker %v", i+1)
		go asyncProcessor(workQueue)

	}
}

func main() {

	// Connect to the RMQ server - @TODO Dynamic worker names, not W1...
	err := conejo.Consume(rmq, queue, exchange, "W1", workQueue)

	// Check to make sure the there were no errors in consuming
	if err != nil {

		// Foobar no wascally wabbits!!
		log.Printf("ERR: Could not consume messages - %q", err)

	} // Consume Messages
}

/* ----------------------------------------------------------------------------
Function used as a goroutine to process the incomming message.
-----------------------------------------------------------------------------*/
func asyncProcessor(ch chan string) {

	// Range over the messages in the channel
	for m := range ch {

		log.Println(m)

	}
}
