package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
        // Kafka broker address
        brokerAddress := "localhost:9092" // Important: use localhost because of the Host listener

        // Topic name
        topic := "raw_data" // Replace with your topic name

        // Create a Kafka writer
        writer := kafka.NewWriter(kafka.WriterConfig{
                Brokers:  []string{brokerAddress},
                Topic:    topic,
                Balancer: &kafka.LeastBytes{}, // Optional: use a balancer
        })
        defer writer.Close()

        // Message to send
        message := kafka.Message{
                Key:   []byte("key1"),
                Value: []byte("Hello, Kafka!"),
                Time:  time.Now(),
        }

        // Send the message
        err := writer.WriteMessages(context.Background(), message)
        if err != nil {
                log.Fatalf("Failed to write message: %v", err)
        }

        fmt.Println("Message sent successfully!")
}
