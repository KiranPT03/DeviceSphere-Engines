package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

// Models
type RawProperty struct {
	PropertyRefID string `json:"referenceId"`
	Value         string `json:"value"`
}

type RawData struct {
	DeviceRefID string        `json:"deviceRefId"`
	Properties  []RawProperty `json:"properties"`
}

// Producer Config
type NatsProducerConfig struct {
	Servers []string
	Subject string
}

// NatsProducer
type NatsProducer struct {
	nc     *nats.Conn
	config NatsProducerConfig
}

// NewNatsProducer
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
	np.nc, err = nats.Connect(np.config.Servers[0])
	if err != nil {
		log.Printf("Error connecting to NATS: %v", err)
		return err
	}

	return nil
}

// Produce
func (np *NatsProducer) Produce(ctx context.Context, key, value []byte) error {
	err := np.nc.Publish(np.config.Subject, value)
	if err != nil {
		log.Printf("Error publishing message to NATS: %v", err)
		return err
	}
	return nil
}

// Close
func (np *NatsProducer) Close() {
	if np.nc != nil {
		np.nc.Drain()
		np.nc.Close()
	}
	log.Println("NATS producer closed.")
}

func main() {
	// Configure the NATS producer
	producerConfig := NatsProducerConfig{
		Servers: []string{"nats://192.168.1.6:4222"}, // Replace with your NATS server address
		Subject: "data.rawData",                      // Replace with your subject
	}

	// Create the NATS producer
	producer, err := NewNatsProducer(producerConfig)
	if err != nil {
		log.Fatalf("Error creating NATS producer: %v", err)
	}
	defer producer.Close()

	app := fiber.New()

	app.Post("kepware/devices/raw-data", func(c *fiber.Ctx) error {
		var rawData RawData
		if err := c.BodyParser(&rawData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		jsonData, err := json.Marshal(rawData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error marshaling JSON",
			})
		}

		ctx := context.Background()
		if err := producer.Produce(ctx, nil, jsonData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error producing message to NATS",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Data sent to NATS",
		})
	})

	err = app.Listen(":4000")
	if err != nil {
		log.Fatalf("Error starting Fiber server: %v", err)
	}
}
