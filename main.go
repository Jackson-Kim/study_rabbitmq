package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	
	amqp "github.com/rabbitmq/amqp091-go"
)

var queueNameA string = "helloA"
var queueNameB string = "helloB"
var exchange string = "myExchange"


func main() {
	e := echo.New()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// controller layer
	e.GET("mq/put/:msg", func(c echo.Context) error {
		msg := c.Param("msg")
		producer(conn, msg)
		return c.String(http.StatusOK, fmt.Sprintf("put message : [%s]", msg))
	})

	err = initExchange(conn)
	failOnError(err, "Failed to initExchange")

	// go customer(conn, queueNameA)
	// go customer(conn, queueNameB)

	e.Logger.Fatal(e.Start(":1323"))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func initExchange(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange,    // exchange.name 
		amqp.ExchangeFanout,    // exchange.kind :: ExchangeDirect, ExchangeFanout, ExchangeTopic, ExchangeHeaders
		false, // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)

	failOnError(err, "Failed to declare exchange")


	
	_, err = ch.QueueDeclare(
		queueNameA, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	_, err = ch.QueueDeclare(
		queueNameB, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(queueNameA, "*", exchange, false, nil)
	failOnError(err, "Failed to binding queue" + queueNameA)
	err = ch.QueueBind(queueNameB, "*", exchange, false, nil)
	failOnError(err, "Failed to binding queue" + queueNameB)
	return err
}

func producer(conn *amqp.Connection, msg string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Publish(
		exchange,     // exchange
		"", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msg)
}

func customer(conn *amqp.Connection, queueName string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		queueName+"_consumer",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message : %s : %s", queueName, d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return
}