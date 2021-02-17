package main

import (
  "log"

  "github.com/streadway/amqp"
)
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}
func main() {
	// establish connection to rabbitmq as a guest user
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connection established!")
	defer conn.Close()
}
