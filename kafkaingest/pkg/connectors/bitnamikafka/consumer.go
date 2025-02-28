package bitnamikafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumerConfig holds the configuration for the KafkaConsumer.
type KafkaConsumerConfig struct {
        Brokers  []string
        GroupID  string
        MinBytes int
        MaxBytes int
        Topics   []string
}

// KafkaConsumer is a struct that encapsulates the Kafka consumer logic.
type KafkaConsumer struct {
        readers  map[string]*kafka.Reader
        msgChans map[string]chan []byte
        config   KafkaConsumerConfig
}

// NewKafkaConsumer creates a new KafkaConsumer instance with the given configuration.
func NewKafkaConsumer(config KafkaConsumerConfig) *KafkaConsumer {
        kc := &KafkaConsumer{
                readers:  make(map[string]*kafka.Reader),
                msgChans: make(map[string]chan []byte),
                config:   config,
        }
        kc.init()
        return kc
}

func (kc *KafkaConsumer) init() {
        for _, topic := range kc.config.Topics {
                reader := kafka.NewReader(kafka.ReaderConfig{
                        Brokers:  kc.config.Brokers,
                        Topic:    topic,
                        GroupID:  kc.config.GroupID,
                        MinBytes: kc.config.MinBytes,
                        MaxBytes: kc.config.MaxBytes,
                })

                kc.readers[topic] = reader
                kc.msgChans[topic] = make(chan []byte)

                go func(topic string) {
                        for {
                                msg, err := kc.readers[topic].ReadMessage(context.Background())
                                if err != nil {
                                        log.Printf("Error reading message from topic %s: %v", topic, err)
                                        //Consider adding retry logic or error handling here.
                                        continue; //or return, depends on how fatal you consider errors.
                                }
                                kc.msgChans[topic] <- msg.Value
                        }
                }(topic)
        }
}

// Consume returns a channel that receives messages from the specified topic.
func (kc *KafkaConsumer) Consume(topic string) <-chan []byte {
        return kc.msgChans[topic]
}

// Close closes all Kafka readers.
func (kc *KafkaConsumer) Close() {
        for _, reader := range kc.readers {
                if err := reader.Close(); err != nil {
                        log.Printf("Error closing Kafka reader: %v", err)
                }
        }
}