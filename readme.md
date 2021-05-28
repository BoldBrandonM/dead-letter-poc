This repo contains a basic rabbitMQ dead-letter-queue implementation in golang leveraging the library github.com/streadway/amqp, and serves as a starting point / testing ground for future endeavors hoping to leverage dead-letter-queues in production applications.

- Commands:
    - Emit messages with topics 'text' and 'bytes' to target queues: `make emit-all`
    - Consume messages with the 'text' topic: `make receive-text`
        - Optionally failing to process messages: `make arg=-fail receive-text`
    - Consume messages with the 'bytes' topic: `make receive-bytes`
        - Optionally failing to process messages: `make arg=-fail receive-bytes`
    - Consume dead-letter-messages, republishing to original exchange: `make receive-dead-letters`
