package Rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

func Receiving() {

	// Initialize Connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Println("Failed Connecting to RabbitMQ")
	}

	defer conn.Close()

	// Initialize Channel
	ch, err := conn.Channel()

	if err != nil {
		log.Println("Failed Opening Channel")
	}

	defer ch.Close()

	// Initialize Queue
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Println("Failed Starting Queue")
	}

	// Receive Data
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Println("Failed Fetching Data")
	}

	forever := make(chan bool)

	go func() {

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}

	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
