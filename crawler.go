package main

import (
	"github.com/streadway/amqp"
	"github.com/vwiart/rabbitmq-client/rabbitmq"
	"github.com/vwiart/rabbitmq-client/config"
	"log"
	"fmt"
)

const (
	CHANNEL_NAME = "logs"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial(config.GetConfig().RabbitMQServerURL)
	handleError(err, rabbitmq.FAILED_TO_CONNECT)
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, rabbitmq.FAILED_TO_OPEN_CHANNEL)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		CHANNEL_NAME,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	handleError(err, rabbitmq.FAILED_TO_DECLARE_EXCHANGE)

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	handleError(err,rabbitmq.FAILED_TO_DECLARE_A_QUEUE )

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		CHANNEL_NAME, // exchange
		false,
		nil)
	handleError(err, rabbitmq.FAILED_TO_BIND_A_QUEUE)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, rabbitmq.FAILED_TO_REGISTER_A_CONSUMER)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
