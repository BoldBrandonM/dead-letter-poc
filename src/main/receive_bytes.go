package main

import (
	"github.com/streadway/amqp"
	"log"
	"dead-letter-poc/helpers"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q := helpers.SetupQueueBinding(ch, "messages", "topic", "bytes_queue", "bytes", "messages_dlx", false)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// TODO: toggle ack vs nack to simulate max retries based on cli flags
			// d.Nack(false, false)

			d.Ack(false)
			log.Printf("Received a bytes message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
