package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)
var cachedResult interface{}
var crMutex sync.RWMutex
func getResponse() (interface{},error) {
    time.Sleep(time.Duration(rand.Int63n(3000) * int64(time.Millisecond)))
    if rand.Intn(4) == 3 {
        return nil, errors.New("third error")
    }
    return time.Now(),nil
}

type responseResult struct {
    res interface{}
    err error
}

func getData() (interface{},error) {    
    deliverResult := make(chan responseResult)
    go func () {
        res, err := getResponse()
        deliverResult <- responseResult{res:res,err:err}
        if err == nil {
            crMutex.Lock()
            cachedResult = res
            crMutex.Unlock()
        }
    }()
    select {
	case rr := <- deliverResult :
        return rr.res,rr.err
	case <- time.After(time.Second) :
        log.Print("got cached result on timeout")
        crMutex.RLock()
        defer crMutex.RUnlock()
        return cachedResult,nil
    }
}

func main() {
    for i := 0; i < 5; i++ {
        var wg sync.WaitGroup
        for j := 0; j < 2; j++ {
            n := i * 2 + j
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