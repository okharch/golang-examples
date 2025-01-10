package main

import (
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from error: %v", err)
		}
	}()

	log.Println("Request received, emulating service call (5 seconds)")

	select {
	case <-r.Context().Done():
		log.Println("Client disconnected before processing completed")
		return
	case <-time.After(5 * time.Second):
		// Simulate service delay
	}

	_, err := w.Write([]byte("Hello, World!\n"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush() // Force immediate flush
	}

	log.Println("Response successfully written to the client")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
