package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	topic := "purchases" // Replace with the Kafka topic you desire

	// Create and send messages continuously
	for {
		message := "Hello, MSIB!" // Message to be sent
		producerMessage := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		}

		// Send the message to Kafka
		partition, offset, err := producer.SendMessage(producerMessage)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		// Display information about the sent message
		fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

		// Wait for 1 second before sending the next message
		time.Sleep(1 * time.Second)
	}
}
