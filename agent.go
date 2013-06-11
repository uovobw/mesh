package main

import (
	//	"fmt"
	"io"
	//	"mesh/target"
	//	"mesh/tests/pingtest"
	"flag"
	"log"
	"mesh/opinion"
	"mesh/tests/pingtest"
	"net/http"
	"runtime"
	"time"
)

// sleep between each tick
const waitTime = 5
const testRefreshTime = 30

var targetList []string
var idleTicker = time.NewTicker(waitTime * time.Second)
var refreshTicker = time.NewTicker(testRefreshTime * time.Second)
var quit = make(chan int)
var debug = flag.Bool("d", false, "Enable debug")
var op = opinion.NewOpinion()

func quitHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		log.Print("Got quit request")
	}
	io.WriteString(w, "quit\n")
	doQuit()
}

func opinionHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		log.Print("Opinion request from ", r.RemoteAddr)
	}
	io.WriteString(w, op.Print())
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		log.Print("Forced by ", r.RemoteAddr)
	}
	go buildOpinion()
	io.WriteString(w, "refresh\n")
}

func doQuit() {
	idleTicker.Stop()
	refreshTicker.Stop()
	quit <- 1
}

func runTest() (err error) {
	return nil
}

func setup() (err error) {
	// placeholder function to fetch data from
	// configurator
	targetList = append(targetList, "127.0.0.1")
	targetList = append(targetList, "192.168.1.1")
	return nil
}

func buildOpinion() {
	log.Print("Building Opinion")
	for _, target := range targetList {
		log.Print(target)
		t := pingtest.Pingtest{}
		t.Setup(map[string]string{"address": target, "count": "2"})
		result, err := t.Run()
		if err != nil {
			log.Print("test failed for target " + target)
		}
		op.SetOpinionForHost(target, result)
	}
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
	http.HandleFunc("/refresh", refreshHandler)
	go http.ListenAndServe(":6060", nil)

	// setup the local Agent
	if setup() != nil {
		panic("Error in setup function")
	}

	for {
		select {
		case <-idleTicker.C:
			if *debug {
				log.Print("Ticked")
			}
		case <-refreshTicker.C:
			if *debug {
				log.Print("Build opinion received")
			}
			go buildOpinion()
		case <-quit:
			log.Print("Quitting")
			return
		}
	}

}
