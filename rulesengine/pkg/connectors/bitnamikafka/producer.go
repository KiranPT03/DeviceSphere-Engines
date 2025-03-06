package bitnamikafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaProducerConfig holds the configuration for the KafkaProducer.
type KafkaProducerConfig struct {
	Brokers []string
	Topic   string
}

// KafkaProducer is a struct that encapsulates the Kafka producer logic.
type KafkaProducer struct {
	writer *kafka.Writer
	config KafkaProducerConfig
}

// NewKafkaProducer creates a new KafkaProducer instance with the given configuration.
func NewKafkaProducer(config KafkaProducerConfig) *KafkaProducer {
	kp := &KafkaProducer{
		config: config,
	}
	kp.init()
	return kp
}

func (kp *KafkaProducer) init() {
	kp.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kp.config.Brokers,
		Topic:    kp.config.Topic,
		Balancer: &kafka.LeastBytes{},
	})
}

// Produce sends a message to the Kafka topic.
func (kp *KafkaProducer) Produce(ctx context.Context, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
	}

	err := kp.writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
		return err
	}
	return nil
}

// Close closes the Kafka writer.
func (kp *KafkaProducer) Close() {
	if err := kp.writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
	}
}
