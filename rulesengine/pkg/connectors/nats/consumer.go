package nats

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

// NatsConsumerConfig holds the configuration for the NatsConsumer.
type NatsConsumerConfig struct {
	Servers  []string
	GroupID  string
	Subjects []string
}

// NatsConsumer is a struct that encapsulates the NATS consumer logic.
type NatsConsumer struct {
	nc       *nats.Conn
	msgChans map[string]chan *nats.Msg
	config   NatsConsumerConfig
	subs     map[string]*nats.Subscription
}

// NewNatsConsumer creates a new NatsConsumer instance with the given configuration.
func NewNatsConsumer(config NatsConsumerConfig) *NatsConsumer {
	nc := &NatsConsumer{
		msgChans: make(map[string]chan *nats.Msg),
		config:   config,
		subs:     make(map[string]*nats.Subscription),
	}
	nc.init()
	return nc
}

func (nc *NatsConsumer) init() {
	var err error
	nc.nc, err = nats.Connect(nc.config.Servers[0]) // Simplification: using first server
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	for _, subject := range nc.config.Subjects {
		nc.msgChans[subject] = make(chan *nats.Msg)

		// Use QueueSubscribe instead of JetStream QueueSubscribe
		sub, err := nc.nc.QueueSubscribe(subject, nc.config.GroupID, func(msg *nats.Msg) {
			nc.msgChans[subject] <- msg
		})
		if err != nil {
			log.Printf("Error subscribing to subject %s: %v", subject, err)
			continue
		}
		nc.subs[subject] = sub
	}

	// Handle graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		nc.Close()
		os.Exit(0)
	}()
}

// Consume returns a channel that receives messages from the specified subject.
func (nc *NatsConsumer) Consume(subject string) <-chan *nats.Msg {
	return nc.msgChans[subject]
}

// Close closes the NATS connection and unsubscribes all subscriptions.
func (nc *NatsConsumer) Close() {
	for _, sub := range nc.subs {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("Error unsubscribing: %v", err)
		}
	}
	if nc.nc != nil {
		nc.nc.Drain()
		nc.nc.Close()
	}

	log.Println("NATS consumer closed.")
}
