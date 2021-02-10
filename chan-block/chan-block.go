package main

import (
	"log"

	"time"
)

func main() {
	ch := make(chan int, 1)
	blocked := true
	go func() {
		ch <- 1
		blocked = false
	}()
	time.Sleep(time.Millisecond * 100)
	if blocked {
		log.Println("Blocked!")
	} else {
		log.Println("Not blocked!")
	}
}
