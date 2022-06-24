package Rabbit

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

func TestSending(t *testing.T) {

	var D RabbitMQData
	D.Inventory = 1000

	var network bytes.Buffer
	err := json.NewEncoder(&network).Encode(D)

	if err != nil {
		log.Println("Error Encoding Object")
	}

	Sending(network.Bytes())

}

func TestReceiving(t *testing.T) {

	Receiving()

}
