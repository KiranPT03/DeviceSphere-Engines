package nats

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

// NatsProducerConfig holds the configuration for the NatsProducer.
type NatsProducerConfig struct {
	Servers []string
	Subject string
}

// NatsProducer is a struct that encapsulates the NATS producer logic.
type NatsProducer struct {
	nc     *nats.Conn
	config NatsProducerConfig
}

// NewNatsProducer creates a new NatsProducer instance with the given configuration.
func NewNatsProducer(config NatsProducerConfig) (*NatsProducer, error) {
	np := &NatsProducer{
		config: config,
	}
	err := np.init()
	if err != nil {
		return nil, err
	}
	return np, nil
}

func (np *NatsProducer) init() error {
	var err error
	np.nc, err = nats.Connect(np.config.Servers[0]) // Simplification: using first server
	if err != nil {
		log.Printf("Error connecting to NATS: %v", err)
		return err
	}
	return nil
}

// Produce sends a message to the NATS subject.
func (np *NatsProducer) Produce(ctx context.Context, key, value []byte) error {
	err := np.nc.Publish(np.config.Subject, value)
	if err != nil {
		log.Printf("Error publishing message to NATS: %v", err)
		return err
	}
	return nil
}

// Close closes the NATS connection.
func (np *NatsProducer) Close() {
	if np.nc != nil {
		np.nc.Drain()
		np.nc.Close()
	}
	log.Println("NATS producer closed.")
}
