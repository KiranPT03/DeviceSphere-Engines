package ruleprocessor

import (
	"encoding/json"

	config "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/config"
	nats "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/connectors/nats"
	postgres "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/databases/postgres"
	models "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/models"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/utils/loggers"
)

type PostgreProcessor struct {
	PGRepository *postgres.PostgreSQLRepository
	NatsConsumer *nats.NatsConsumer
	NatsProducer *nats.NatsProducer
	Config       *config.Config
}

func NewPostgreProcessor(config *config.Config) *PostgreProcessor {
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

	return &PostgreProcessor{
		PGRepository: repository,
		NatsConsumer: consumer,
		NatsProducer: producer,
		Config:       config,
	}
}

func (pgp *PostgreProcessor) dataTransformer(data string) {
	log.Debug("Processing received data: %s", data)

	var deviceRec models.Device
	err := json.Unmarshal([]byte(data), &deviceRec)
	if err != nil {
		log.Error("Error unmarshalling JSON: %v", err)
		return
	}

	log.Debug("Device model %+v\n", deviceRec)
	for _, prop := range deviceRec.Properties {
		// Check property ID exist
		log.Debug("Property ID: %s", prop.ID)
		dataToUpdate := map[string]interface{}{
			"value": prop.Value,
		}
		updateResult, updateErr := pgp.PGRepository.Update("properties", prop.ID, dataToUpdate)
		if updateErr != nil {
			log.Error("Error while checking existance, %v", updateErr)
		}
		log.Debug("Update result: %d", updateResult)

	}

}

func (pgp *PostgreProcessor) ProcessData() {
	log.Info("Processing the data")
	consumer := pgp.NatsConsumer
	defer consumer.Close()

	producer := pgp.NatsProducer
	defer producer.Close()

	rawDataTopicChan := consumer.Consume(pgp.Config.Nats.InletSubject)
	// processedDataTopicChan := consumer.Consume("other_subject")

	for {
		select {
		case msg := <-rawDataTopicChan:
			go pgp.dataTransformer(string(msg.Data))
			// case msg := <-processedDataTopicChan:
			// 	log.Debug("Received message topic: processed_data: %s", msg.Data)
		}
	}
}
