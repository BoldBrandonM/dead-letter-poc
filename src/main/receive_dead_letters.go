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

	msgs, err := ch.Consume(
		"messages_dlq", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		numProcessed := 0
		for d := range msgs {
			log.Printf("Received a text message: %s", d.Body)
			log.Printf("exchange: %v\nrouting key: %v", d.Exchange, d.RoutingKey)
			err = ch.Publish(
				"messages",     // exchange
				d.RoutingKey, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing {
					ContentType: "application",
					Body:        d.Body,
				})
			helpers.FailOnError(err, "Failed to publish a message")
			log.Printf(" [x] Resurrected message %v: %s", numProcessed, d.Body)

			// sleep used here to limit spam in case of unprocessable messages
			time.Sleep(1 * time.Second)
			numProcessed++
		}
	}()

	log.Printf(" [*] Waiting for dead letters. To exit press CTRL+C")
	<-forever
}
