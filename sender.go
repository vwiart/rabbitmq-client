package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
	"github.com/vwiart/rabbitmq-client/rabbitmq"
	"github.com/vwiart/rabbitmq-client/config"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial(config.GetConfig().RabbitMQServerURL)
	failOnError(err, rabbitmq.FAILED_TO_CONNECT)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, rabbitmq.FAILED_TO_OPEN_CHANNEL)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, rabbitmq.FAILED_TO_DECLARE_EXCHANGE)

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, rabbitmq.FAILED_TO_PUBLISH_MESSAGE)

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}