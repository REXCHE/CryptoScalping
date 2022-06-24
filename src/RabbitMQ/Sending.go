package Rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

/*
	This Method Sends Data to a Port Using RabbitMQ

	Input:

	1. Object Converted to Bytes

*/
func Sending(byte_arr []byte) {

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

	// Publish Data
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byte_arr,
		})

	if err != nil {
		log.Println("Failed Publishing Data")
	}

}
