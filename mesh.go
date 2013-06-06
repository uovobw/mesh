package main

import (
	"fmt"
	"io"
	"mesh/target"
	"mesh/tests/pingtest"
	"net/http"
	"time"
)

var msgChan = make(chan string, 1)

func MainRoute(w http.ResponseWriter, r *http.Request) {
	msgChan <- r.RequestURI
	w.WriteHeader(200)
	io.WriteString(w, "ok\n")
}

func DoPing(w http.ResponseWriter, r *http.Request) {
	RunTest("127.0.0.1", 3)
}

func RunTest(address string, count int) {
	pingTest := pingtest.MakeNewPingTest(address, count)
	_, err := pingTest.Run()
	if err == nil {
		fmt.Println("Success")
	} else {
		fmt.Println("Failed with error:", err)
	}
}

func main() {

	http.HandleFunc("/", MainRoute)
	http.HandleFunc("/ping", DoPing)

	t := target.MakeNewTarget("testtarget", "1.2.3.4")
	target.Print()
	t.Print()

	server := &http.Server{
		Addr:           ":8080",
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Server created")

	go server.ListenAndServe()

	for {
		fmt.Println("Looping")
		requestURI := <-msgChan
		fmt.Printf("Got request for %s\n", requestURI)
		if requestURI == "/quit" {
			fmt.Println("Exiting")
			break
		}
	}
}
