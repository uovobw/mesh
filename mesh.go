package main

import (
    "net/http"
    "fmt"
    "time"
    "io"
    "mesh/target"
    "mesh/tests/pingtest"
    )

var msgChan = make(chan string, 1)

func MainRoute(w http.ResponseWriter, r *http.Request) {
    msgChan <- r.RequestURI
    w.WriteHeader(200)
    io.WriteString(w, "ok\n")
}

func RunTest() {
    pingTest := pingtest.MakeNewPingTest("127.0.0.1", 3)
    _, err := pingTest.Run()
    if err == nil {
        fmt.Println("Success")
    } else {
        fmt.Println("Failed with error:", err)
    }
}



func main() {

    http.HandleFunc("/", MainRoute)


    t := target.MakeNewTarget("chiappette", "1.2.3.4")
    target.Print()
    t.Print()

    server := &http.Server {
        Addr:           ":8080",
        Handler:        http.DefaultServeMux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    fmt.Println("Server created")

    go server.ListenAndServe()

    RunTest()

    for {
        fmt.Println("Looping")
        requestURI := <-msgChan
        fmt.Printf("Got request for %s\n", requestURI)
        if requestURI == "/chiappe" {
            fmt.Println("Exiting")
            break
        }
    }
}
