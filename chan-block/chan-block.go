package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// depending on bufSize of channel writing to it either blocks or not
	// if its size is 0 it blocks, if it is at least one it does not
	// this program demonstrates that
	defer log.Println("Exiting from program")
	log.Println("Program started")
	var bufSize int
	if len(os.Args) == 2 {
		s, err := strconv.Atoi(os.Args[1])
		if err != nil {
			s = 0
		}
		bufSize = s
	}
	ch := make(chan int, bufSize)
	notBlocked := make(chan bool)
	go func() {
		ch <- 1
		notBlocked <- true
	}()
	// looping while blocked in go routine
	for {
		select {
		case <-time.After(time.Second):
			log.Println("Blocked! Unblocking...")
			<-ch
		case <-notBlocked:
			log.Println("Not blocked!")
			return
		}
	}

}
