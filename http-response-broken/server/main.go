package main

import (
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")

	// Simulate processing delay
	select {
	case <-time.After(1 * time.Second): // First pause for simplicity
	case <-r.Context().Done():
		log.Println("Client disconnected before processing completed")
		return
	}

	// First write and flush
	_, err := w.Write([]byte("First line of the response\n"))
	if err != nil {
		log.Printf("Error writing first part of the response: %v", err)
		return
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush() // Flush after the first write
	}
	log.Println("Flushed first line of the response")

	select {
	case <-time.After(1 * time.Second): // First pause for simplicity
	case <-r.Context().Done():
		log.Println("Client disconnected before second line completed")
		return
	}

	// Second write and flush
	_, err = w.Write([]byte("Second line of the response\n"))
	if err != nil {
		log.Printf("Error writing second part of the response: %v", err)
		return
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush() // Flush after the second write
	}
	// how we can detect that Write was unsuccessful, flush does not have any error status?

	log.Println("Flushed second line of the response, waiting for 1 second")
	time.Sleep(1 * time.Second)
	select {
	case <-r.Context().Done():
		log.Println("Client disconnected before handler finished")
		return
	default:
		log.Println("All the lines were sent")
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
