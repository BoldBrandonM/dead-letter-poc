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

	err = ch.ExchangeDeclare(
		"messages", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	helpers.FailOnError(err, "Failed to declare an exchange")

	body := "Hello World!"
	err = ch.Publish(
		"messages",     // exchange
		"text", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	helpers.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	err = ch.Publish(
		"messages",     // exchange
		"bytes", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "application",
			Body:        []byte(body),
		})
	helpers.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent bytes for: %s", body)


	err = ch.Publish(
		"messages",     // exchange
		"text", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	helpers.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
