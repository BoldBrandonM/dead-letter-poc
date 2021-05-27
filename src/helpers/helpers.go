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

func SetupQueue(ch *amqp.Channel, exchangeName, exchangeKind, deadLetterExchangeName, queueName, routingKey string) (amqp.Queue, error) {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		exchangeKind,      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

	textQ, err := ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange": deadLetterExchangeName},     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(textQ.Name,       // queue name
		routingKey,            // routing key
		exchangeName, // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind a queue")

	if err != nil {
		return amqp.Queue{}, err
	}

	return textQ, nil
}
