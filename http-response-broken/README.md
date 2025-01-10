# Demo: REST Server and Client with Connection Handling
This project demonstrates how a REST server can detect broken client connections during the request handling and response writing phases. The demo uses a server-client pair to orchestrate various scenarios, such as waiting for a full response or simulating premature disconnections at specific points in the server's lifecycle.
## Overview
### Server
The server implements a REST handler (`http.HandlerFunc`) that:
1. Simulates request processing delays.
2. Sends a two-part streaming response with a pause and flush between writes.
3. Detects and logs when a client disconnects at various stages using the Go `context.Context` cancellation mechanism.

### Client
The client interacts with the server and:
1. Simulates successful requests that receive the full response.
2. Simulates premature disconnections at different stages:
   - Before the server starts writing the response.
   - After receiving the first flush (first part of the response).
   - After the full response but before the server completes the handler.

## Features
### Server
- Handles HTTP requests with simulated delays to represent realistic processing.
- Splits the response into two writes:
   - The first part is flushed immediately after a 1-second delay.
   - The second part is written and flushed after another 1-second delay.

- Detects and logs client disconnections during all phases:
   - Pre-response processing.
   - Between the two writes.
   - After the final write but before the handler fully completes.

### Client
- Waits for the server to send a full response and logs success for normal requests.
- Disconnects mid-request to simulate the following scenarios:
   - Before the server begins writing the response.
   - Before the server sends the second part of the response.
   - After the full response but before the server's handler completes.

## How It Works
### Server Implementation
The server:
1. **Simulates Initial Processing Delay:** Uses a `time.Sleep(1 * time.Second)` before beginning the response. A `context.Done` check ensures that if the client disconnects during this delay, it is logged as `Client disconnected before processing completed`.
2. **Writes and Flushes in Parts:**
   - **First Write:** Sends the first line of the response with `Flush()` after completing the initial delay. Detects disconnections before or during the flush.
   - **Second Write:** Sends the second line with a 1-second pause before writing. Detects disconnections in between.

3. **Final Wait:** Introduces a final 1-second delay after the second flush and detects if the client disconnects after the full response but before the handler finishes.

### Client Implementation
The client:
1. **Full Response Test:** The client waits for the server to send the complete response without disconnecting.
2. **Premature Disconnections:** The client closes the connection deliberately under specific timing scenarios:
   - Before the server has a chance to respond.
   - After the first line of the response is received and flushed.
   - After receiving the full response but before the server’s handler finishes cleanup.

### Logs and Results
- **Logs on Successful Response:** The server logs completion of all response-handling stages.
- **Logs on Disconnections:** The server logs specific messages indicating where the disconnection occurred:
   - Before processing: `Client disconnected before processing completed`
   - Before second write: `Client disconnected before second line completed`
   - Before final handler finish: `Client disconnected before handler finished`

## Running the Demo
### Prerequisites
- Go 1.23 or later installed.
- Network access (e.g., localhost).

### Steps to Run
#### 1. Start the Server
Run the server:
``` bash
go run server/main.go
```
The server starts on port `:8080` and logs connection statuses.
#### 2. Run the Client
Run the client to simulate all scenarios:
``` bash
go run client/main.go
```
The client runs the following test cases and logs their results:
1. **Full Response Wait:** The client waits for the complete response and logs both lines as received.
2. **Premature Disconnection #1:** The client sends a request and disconnects after 0.5 seconds.
3. **Premature Disconnection #2:** The client disconnects after 1.5 seconds, after receiving the first response part.
4. **Premature Disconnection #3:** The client disconnects after the full response but before the server finishes cleanup.

### Expected Logs
#### Server Logs
For each client test, the server logs messages showing where the client disconnected or if the request completed successfully.
Example:
1. **Successful Request:**
``` plaintext
   Request received
   Flushed first line of the response
   Flushed second line of the response
   All the lines were sent
```
1. **Disconnect Before Processing:**
``` plaintext
   Request received
   Client disconnected before processing completed
```
1. **Disconnect Before Second Line:**
``` plaintext
   Request received
   Flushed first line of the response
   Client disconnected before second line completed
```
1. **Disconnect After Full Response:**
``` plaintext
   Request received
   Flushed first line of the response
   Flushed second line of the response
   Client disconnected before handler finished
```
#### Client Logs
The client logs each scenario appropriately. Example:
1. **Successful Request:**
``` plaintext
   Case 1: Wait for full response
   Response line: First line of the response
   Response line: Second line of the response
   Case 1: Request completed successfully
```
1. **Premature Disconnection:**
``` plaintext
   Case 2: Disconnect after 0.5 seconds
   Disconnecting from the server
   Case 2: Disconnect completed
```
## Key Learnings
This demo illustrates:
- How to handle broken client connections gracefully in a REST server.
- How to use response flushing (`http.Flusher`) to test partial responses for real-time client-server communication.
- How to use `context.Context` to detect client disconnections during different stages of request processing.

## Future Improvements
- Add metrics to track disconnection frequencies and patterns.
- Extend the client to include custom headers or payloads for more realistic scenarios.
- Enhance server error detection mechanisms, such as handling errors explicitly on write or flush operations.
