package helpers

import (
	"log"
	"github.com/streadway/amqp"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SetupQueueBinding(ch *amqp.Channel, exchangeName, exchangeKind, queueName, routingKey, deadLetterExchangeName string, durableQueue bool) amqp.Queue {
	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare an exchange")

	var qArgs amqp.Table
	if deadLetterExchangeName != "" {
		qArgs = amqp.Table{"x-dead-letter-exchange": deadLetterExchangeName}
	}

	q, err := ch.QueueDeclare(
		queueName,
		durableQueue,
		false,
		false,
		false,
		qArgs,
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	FailOnError(err, "Failed to bind a queue")

	return q
}

func PublishMessage(ch *amqp.Channel, exchangeName, routingKey, contentType, body string) {
	err := ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing {
			ContentType: contentType,
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s for %s", routingKey, body)
}
