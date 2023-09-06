package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gin-kafka-ws/websocket"
)

func main() {
	// Initialize the Gin router
	r := gin.Default()

	// Set up the WebSocket endpoint
	r.GET("/ws", websocket.HandleWebSocket)

	// Start the WebSocket server asynchronously
	go websocket.StartWebSocketServer()

	// Kafka consumer configuration
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	// Create a Kafka partition consumer for the "purchases" topic
	partitionConsumer, err := consumer.ConsumePartition("purchases", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Run the Gin HTTP server on port 9999
	r.Run(":9999")

	num := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Received message %d: %s\n", num, string(msg.Value))
			num++

			// Send the message to all connected WebSocket clients
			websocket.BroadcastMessage(string(msg.Value))

			// Send a WebSocket update with a notification message
			websocket.SendWebSocketUpdate("This is a notification")

		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error: %v\n", err.Err)
		}
	}
}