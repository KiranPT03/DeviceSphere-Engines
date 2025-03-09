package notificationprocessor

import (
	config "factorysphere/devicesphere/engines/notificationengine/pkg/config"
	nats "factorysphere/devicesphere/engines/notificationengine/pkg/connectors/nats"
	postgres "factorysphere/devicesphere/engines/notificationengine/pkg/databases/postgres"
	log "factorysphere/devicesphere/engines/notificationengine/pkg/utils/loggers"
)

type NotificationProcessor struct {
	PGRepository *postgres.PostgreSQLRepository
	NatsConsumer *nats.NatsConsumer
	NatsProducer *nats.NatsProducer
	Config       *config.Config
}

func NewNotificationProcessor(config *config.Config) *NotificationProcessor {
	// Create a new PostgreSQL repository instance
	repository, err := postgres.NewPostgreSQLRepository(config)
	if err != nil {
		panic(err)
	}

	natsConsumerConfig := nats.NatsConsumerConfig{
		Servers:  []string{config.Nats.Server},
		Stream:   config.Nats.Stream,
		GroupID:  config.Nats.ConsumerGroup,
		Subjects: []string{config.Nats.InletSubject},
		Durable:  config.Nats.Durable,
	}
	consumer := nats.NewNatsConsumer(natsConsumerConfig)

	natsProducerConfig := nats.NatsProducerConfig{
		Servers: []string{config.Nats.Server},
		Stream:  config.Nats.Stream,
		Subject: config.Nats.OutletSubject,
	}

	producer, err := nats.NewNatsProducer(natsProducerConfig)
	if err != nil {
		panic(err)
	}

	return &NotificationProcessor{
		PGRepository: repository,
		NatsConsumer: consumer,
		NatsProducer: producer,
		Config:       config,
	}
}


func (np *NotificationProcessor) dataTransformer(data string) {
	log.Debug("Data received for rule %s",data)

}

func (np *NotificationProcessor) ProcessData() {
	log.Info("Processing the data")
	consumer := np.NatsConsumer
	defer consumer.Close()

	producer := np.NatsProducer
	defer producer.Close()

	rawDataTopicChan := consumer.Consume(np.Config.Nats.InletSubject)
	// processedDataTopicChan := consumer.Consume("test.topic")

	for {
		select {
		case msg := <-rawDataTopicChan:
			go np.dataTransformer(string(msg.Data))
		// case msg := <-processedDataTopicChan:
		// 	log.Debug("Received message topic: processed_data: %s", msg.Data)
		}
	}
}
