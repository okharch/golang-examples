package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	// Case 1: Successful request, wait for full response
	runClientTest("Case 1: Wait for full response", 0, true)

	// Case 2: Disconnect after 0.5 seconds
	runClientTest("Case 2: Disconnect after 0.5 seconds", 500*time.Millisecond, false)

	// Case 3: Disconnect after 1.5 seconds
	runClientTest("Case 3: Disconnect after 1.5 seconds", 1500*time.Millisecond, false)

	// Case 4: Disconnect after 2.5 seconds
	runClientTest("Case 4: Disconnect after 2.5 seconds", 2500*time.Millisecond, false)
}

func runClientTest(caseDescription string, disconnectAfter time.Duration, waitForFullResponse bool) {
	log.Println(caseDescription)

	if waitForFullResponse {
		// Full response case
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Read response from the server
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			log.Printf("Response line: %s", scanner.Text())
		}

		// Check if there was an error reading in-progress
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading response: %v", err)
		} else {
			log.Println("Response received successfully")
		}
		return
	}

	// Disconnect cases
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")

	// Wait for a specified time before breaking the connection
	time.Sleep(disconnectAfter)

	// Report and close the connection
	log.Println("Disconnecting from the server")
	conn.Close()

	// Wait to give the server time to detect disconnection
	time.Sleep(2 * time.Second)
	log.Println(caseDescription + " completed")
}
