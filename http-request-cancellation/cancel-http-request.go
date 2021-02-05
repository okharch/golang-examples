package main

import (
	"log"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
	"flag"
)

var stepCount int
var stepLock sync.Mutex

func logStep(args ...interface{}) {
	stepLock.Lock()
	defer stepLock.Unlock()
	stepCount++
	//args.unshift
	args = append([]interface{}{fmt.Sprintf("%02d.",stepCount)}, args...)
	log.Println(args...)
}

func main() {
	requestCancelDelay  := flag.Duration("request-timeout",2 * time.Second ,"timeout for expiring request")
	serverResponseDelay  := flag.Duration("server-delay",3 * time.Second ,"programmatic delay for debug server")
	flag.Parse()
	httpServerDone := make(chan struct{})
	gracefulExit := make(chan struct{})
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		duration := *serverResponseDelay
		logStep("Server response delay is ", duration)
		select {
		case <-gracefulExit:
			logStep("Graceful exit from http server")
		case <-time.After(duration):
			logStep(duration, " response successfuly finished!")
		}
		httpServerDone <- struct{}{}
	}))
	defer func() {
		svr.Close()
		logStep("server resources closed!")
	}()
	logStep("making request", svr.URL)
	tr := &http.Transport{} // TODO: copy defaults from http.DefaultTransport
	req, _ := http.NewRequest("GET", svr.URL, nil)
	c := make(chan error, 1)
	go func() {
		logStep("making long request in go routine...")
		client := &http.Client{Transport: tr}
		_, err := client.Do(req)
		// handle response ...
		c <- err
	}()

	// Simulating user cancel request channel
	duration := *requestCancelDelay
	serverDone := false
	logStep("waiting ", duration, " before cancelling request...")
	select {
	case <-httpServerDone:
		logStep("http server done")
		serverDone = true
	case <-time.After(duration):
		logStep("Cancelling request")
		tr.CancelRequest(req)
	}
	logStep("Finding out response error if any (nil means no error) ...")
	err := <-c
	logStep("Response status: ", err)
	if !serverDone {
		logStep("Sending to server response handler request for graceful exit...")
		gracefulExit <- struct{}{}
		logStep("Waiting for graceful exit from server response handler")
		<-httpServerDone
	}
	logStep("End of main")
}
