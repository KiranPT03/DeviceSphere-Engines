package tests

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func TestKafka() {
	// Create a new Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "test_topic",
		Balancer: &kafka.LeastBytes{},
	})

	// Produce messages
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Message test %d", i)
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: []byte(msg),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Produced message: %s\n", msg)
	}

	// Close the writer
	err := writer.Close()
	if err != nil {
		log.Fatal(err)
	}
}
