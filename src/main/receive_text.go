package main

import (
	"github.com/streadway/amqp"
	"log"
	"dead-letter-poc/helpers"
	"flag"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q := helpers.SetupQueueBinding(ch, "messages", "topic", "text_queue", "text", "messages_dlx", false)
	helpers.SetupQueueBinding(ch, "messages_dlx", "fanout", "messages_dlq", "", "", true)

	failFlag := flag.Bool("fail", false, "-fail")
	flag.Parse()
	shouldFail := *failFlag

	// NOTE: Consume is destructive, if set up to auto ack, the messages will be lost between when
	// they are taken from the queue and when they are processed here IF the consumer crashes for whatever reason.
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
			if shouldFail {
				d.Nack(false, false)
			} else {
				d.Ack(false)
				log.Printf("Received a text message: %s", d.Body)

				// NOTE: this demonstrates that we can manually dead-letter a message, for example if we ran out of retries
				if string(d.Body) == "Hello Underworld!" {
					log.Print("Manually publishing dead letter")

					helpers.PublishMessage(ch, "messages_dlx", d.RoutingKey, d.ContentType, string(d.Body))
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
