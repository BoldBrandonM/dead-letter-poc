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

	setupQueues(ch)

	body := "Hello World!"
	log.Printf(" [*] Emitting messages. To exit press CTRL+C")
	for {
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

		time.Sleep(1 * time.Second)
	}
}

func setupQueues(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"messages", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	helpers.FailOnError(err, "Failed to declare an exchange")

	textQ, err := ch.QueueDeclare(
		"text_queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange": "messages_dlx"},     // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(textQ.Name,       // queue name
		"text",            // routing key
		"messages", // exchange
		false,
		nil)
	helpers.FailOnError(err, "Failed to bind a queue")

	bytesQ, err := ch.QueueDeclare(
		"bytes_queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange": "messages_dlx"},     // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(bytesQ.Name,       // queue name
		"bytes",            // routing key
		"messages", // exchange
		false,
		nil)
	helpers.FailOnError(err, "Failed to bind a queue")

	err = ch.ExchangeDeclare(
		"messages_dlx",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	helpers.FailOnError(err, "Failed to declare a dead-letter exchange")

	dlQ, err := ch.QueueDeclare(
		"messages_dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	helpers.FailOnError(err, "Failed to declare a dead-letter queue")

	err = ch.QueueBind(
		dlQ.Name,
		"",
		"messages_dlx",
		false,
		nil,
	)
	helpers.FailOnError(err, "Failed to bind a dead-letter queue")
}
