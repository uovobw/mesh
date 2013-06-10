package main

import (
	//	"fmt"
	"io"
	//	"mesh/target"
	//	"mesh/tests/pingtest"
	"flag"
	"log"
	"mesh/opinion"
	"net/http"
	"runtime"
	"time"
)

// sleep between each tick
const waitTime = 5

const (
	BUILD_OPINION = iota
	TICK
)

var commandChan = make(chan int, 10)
var ticker = time.NewTicker(waitTime * time.Second).C
var debug = flag.Bool("d", false, "Enable debug")
var op = opinion.NewOpinion()

func quitHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		log.Print("Got quit request")
	}
	io.WriteString(w, "quit")
	// close the command chan to trigger exit in the main thread
	close(commandChan)
}

func opinionHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		log.Print("Opinion request from ", r.RemoteAddr)
	}
	io.WriteString(w, op.Print())
}

func runTest() (err error) {
	return nil
}

func buildOpinion() {
	log.Print("Building Opinion")
	op.SetOpinionForHost("aHost", 0.234)
}

func main() {
	// parse cmdline parameters
	flag.Parse()
	log.Print("Starting up")
	// parallelize execution
	runtime.GOMAXPROCS(runtime.NumCPU())

	// create the http server and related routes
	http.HandleFunc("/quit", quitHandler)
	http.HandleFunc("/opinion", opinionHandler)
	go http.ListenAndServe(":6060", nil)

	// build the first opinion
	commandChan <- BUILD_OPINION

	// add a tick each time the ticker ticks
	go func() {
		for {
			<-ticker
			commandChan <- TICK
		}
	}()

	for cmd := range commandChan {
		switch cmd {
		case BUILD_OPINION:
			if *debug {
				log.Print("Build opinion received")
			}
			buildOpinion()
		case TICK:
			if *debug {
				log.Print("Ticked")
			}
		}
	}
	log.Print("Quitting")
}
