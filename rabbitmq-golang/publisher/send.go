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
