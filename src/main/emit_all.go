package main

import (
	"github.com/streadway/amqp"
	"log"
	"dead-letter-poc/helpers"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	helpers.SetupQueueBinding(ch, "messages", "topic", "text_queue", "text", "messages_dlx", false)
	helpers.SetupQueueBinding(ch, "messages", "topic", "bytes_queue", "bytes", "messages_dlx", false)
	helpers.SetupQueueBinding(ch, "messages_dlx", "fanout", "messages_dlq", "", "", true)

	log.Printf(" [*] Emitting messages. To exit press CTRL+C")

	i := 0
	for {
		body := "Hello World!"

		// every 5th message, emit something different. Consumers can use this to test manually publishing dead letters
		if i % 10 == 0 {
			body = "Hello Underworld!"
		}

		helpers.PublishMessage(ch, "messages", "text", "text/plain", body)
		helpers.PublishMessage(ch, "messages", "bytes", "application", body)
		i++

		// sleep used here to simulate metered production of messages
		time.Sleep(1 * time.Second)
	}
}
