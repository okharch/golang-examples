package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var bSize int
	if len(os.Args) == 2 {
		s, err := strconv.Atoi(os.Args[1])
		if err != nil {
			s = 0
		}
		bSize = s
	}
	ch := make(chan int, bSize)
	notBlocked := make(chan bool)
	go func() {
		ch <- 1
		notBlocked <- true
	}()
	select {
	case <-time.After(time.Second):
		log.Println("Blocked!")
	case <-notBlocked:
		log.Println("Not blocked!")
	}

}
