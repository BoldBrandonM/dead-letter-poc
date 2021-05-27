package main

import (
	"dead-letter-poc/helpers"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// NOTE: Consume is destructive, if set up to auto ack, the messages will be lost between when
	// they are taken from the queue and when they are processed here IF the consumer crashes for whatever reason.
	msgs, err := ch.Consume(
		"messages_dlq",
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
		numProcessed := 0
		for d := range msgs {
			log.Printf("exchange: %v\nrouting key: %v\ncontent-type: %v", d.Exchange, d.RoutingKey, d.ContentType)
			err = ch.Publish(
				"messages",
				d.RoutingKey,
				false,
				false,
				amqp.Publishing {
					ContentType: d.ContentType,
					Body:        d.Body,
				})
			helpers.FailOnError(err, "Failed to publish a message")
			log.Printf(" [x] Resurrected message %v: %s", numProcessed, d.Body)
			d.Ack(false)
			numProcessed++

			// sleep used here to limit spam in case of unprocessable messages
			time.Sleep(1 * time.Second)
		}
	}()

	log.Printf(" [*] Waiting for dead letters. To exit press CTRL+C")
	<-forever
}
