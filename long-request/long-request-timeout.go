package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

func getResponse() (interface{}, error) {
	time.Sleep(time.Duration(rand.Int63n(3000) * int64(time.Millisecond)))
	if rand.Intn(4) == 3 {
		return nil, errors.New("third error")
	}
	return time.Now(), nil
}

type responseResult struct {
	res interface{}
	err error
}

// getData make request to getResponse func but waits no more than 1 second(timeout).
// when timed out it returns cachedResult
var cachedResult interface{}
var crMutex sync.RWMutex

// crMutex protects shared variable cachedResult against race condition during read/write operations
func getData() (interface{}, error) {
	// use deliverResult channel to communicate result from goroutine closure into getData
	deliverResult := make(chan responseResult)
	go func() {
		res, err := getResponse()
		deliverResult <- responseResult{res: res, err: err}
		if err == nil {
			// if successful update cachedResult
			crMutex.Lock()
			cachedResult = res
			crMutex.Unlock()
		}
	}()
	// use select to either obtain result from deliverResult channel or exit by timeout
	select {
	case rr := <-deliverResult:
		return rr.res, rr.err
	case <-time.After(time.Second):
		log.Print("got cached result on timeout")
		// use RLock to protect against race condition on shared cachedResult variable
		crMutex.RLock()
		defer crMutex.RUnlock()
		return cachedResult, nil
	}
}

func main() {
	// make outer loop to test 5 internal series of parallel requests
	for i := 0; i < 5; i++ {
		// do requests to getData in Parallel to emulate "real world" scenario
		var wg sync.WaitGroup
		for j := 0; j < 5; j++ {
			n := i*2 + j
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := getData()
				if err != nil {
					log.Println(n, "error:", err)
				} else {
					log.Println(n, res)
				}
			}()
		}
		wg.Wait()
	}
}
