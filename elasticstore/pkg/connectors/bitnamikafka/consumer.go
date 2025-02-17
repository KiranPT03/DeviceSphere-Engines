package bitnamikafka

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	instance *KafkaConsumer
	once     sync.Once
)

type KafkaConsumer struct {
	readers  map[string]*kafka.Reader
	msgChans map[string]chan []byte
}

func GetKafkaConsumer() *KafkaConsumer {
	once.Do(func() {
		instance = &KafkaConsumer{
			readers:  make(map[string]*kafka.Reader),
			msgChans: make(map[string]chan []byte),
		}
		instance.init()
	})
	return instance
}

func (kc *KafkaConsumer) init() {
	topics := []string{"my_topic", "test_topic"} // Append topic here to read on multiple topics

	for _, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			GroupID:  "my_group",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})

		kc.readers[topic] = reader
		kc.msgChans[topic] = make(chan []byte)

		go func(topic string) {
			for {
				msg, err := kc.readers[topic].ReadMessage(context.Background())
				if err != nil {
					log.Fatal(err)
				}
				kc.msgChans[topic] <- msg.Value
			}
		}(topic)
	}
}

func (kc *KafkaConsumer) Consume(topic string) <-chan []byte {
	return kc.msgChans[topic]
}
