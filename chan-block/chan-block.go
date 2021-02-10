package main

import (
	"log"
	"time"
	"os"
	"strconv"
)

func main() {
	var bSize int
	if (len(os.Args) == 2) {
		s, err := strconv.Atoi(os.Args[1])
		if err != nil {
			s = 0
		}
		bSize = s
	}
	ch := make(chan int, bSize)
	blocked := make(chan bool)
	go func() {
		ch <- 1
		blocked <- false
	}()
	select {
	case <-time.After(time.Millisecond * 100):
		log.Println("Blocked!")
	case <-blocked:
		log.Println("Not blocked!")
	}

}
