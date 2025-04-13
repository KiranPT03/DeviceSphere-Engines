package notificationprocessor

import (
	"encoding/json"

	config "factorysphere/devicesphere/engines/notificationengine/pkg/config"
	nats "factorysphere/devicesphere/engines/notificationengine/pkg/connectors/nats"
	elastic "factorysphere/devicesphere/engines/notificationengine/pkg/databases/elastic"
	postgres "factorysphere/devicesphere/engines/notificationengine/pkg/databases/postgres"
	models "factorysphere/devicesphere/engines/notificationengine/pkg/models"
	commons "factorysphere/devicesphere/engines/notificationengine/pkg/utils/commons"
	log "factorysphere/devicesphere/engines/notificationengine/pkg/utils/loggers"

	"github.com/google/uuid"
)

type NotificationProcessor struct {
	PGRepository *postgres.PostgreSQLRepository
	ESRepository *elastic.ElasticsearchRepository
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

	// Create a new Elasticsearch repository instance
	esRepository, err := elastic.NewElasticsearchRepository(config)
	if err != nil {
		log.Error("Failed to create Elasticsearch repository: %v", err)
		panic(err)
	}

	natsConsumerConfig := nats.NatsConsumerConfig{
		Servers: []string{config.Nats.Server},

		GroupID:  config.Nats.ConsumerGroup,
		Subjects: []string{config.Nats.InletSubject},
	}
	consumer := nats.NewNatsConsumer(natsConsumerConfig)

	natsProducerConfig := nats.NatsProducerConfig{
		Servers: []string{config.Nats.Server},
		Subject: config.Nats.OutletSubject,
	}

	producer, err := nats.NewNatsProducer(natsProducerConfig)
	if err != nil {
		panic(err)
	}

	return &NotificationProcessor{
		PGRepository: repository,
		ESRepository: esRepository,
		NatsConsumer: consumer,
		NatsProducer: producer,
		Config:       config,
	}
}

func (np *NotificationProcessor) ruleDataTransformer(data string) {
	log.Debug("Raw rule data received: %s", data)

	// Create a Rule instance to hold the parsed data
	var rule models.Rule

	// Unmarshal the JSON string into the Rule struct
	err := json.Unmarshal([]byte(data), &rule)
	if err != nil {
		log.Error("Failed to unmarshal rule data: %v", err)
		return
	}

	// Generate a random UUID for the notification ID
	notificationID := uuid.New().String()
	currentTime := commons.GetCurrentUTCTimestamp()

	// Create a notification from the rule data
	notification := models.Notification{
		ID:          notificationID,
		Type:        "Rule",
		Severity:    rule.Severity,
		ReceivedAt:  currentTime,
		GeneratedAt: rule.GeneratedAt, // Using CreatedAt as GeneratedAt
		Source:      "RuleEngine",
		Name:        rule.Name,
		AlertID:     rule.ID,
		Description: rule.Description, // Using Name as Description since Rule doesn't have a Description field
		Status:      "Broadcasted",
		BroadcastedStateMetadata: models.BroadcastedStateMetadata{
			BroadcastedBy: "System",
			BroadcastedAt: currentTime,
		},
	}

	// Convert notification to JSON
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		log.Error("Failed to marshal notification: %v", err)
		return
	}

	// Log the created notification
	log.Debug("Created notification: %s", string(notificationJSON))

	// // Send the notification to the outlet topic
	// err = np.NatsProducer.Produce(context.Background(), nil, notificationJSON)
	// if err != nil {
	// 	log.Error("Failed to produce notification message: %v", err)
	// 	return
	// }

	// log.Info("Notification sent to outlet topic: %s", notification.ID)

	// Store the notification in Elasticsearch if available
	if np.ESRepository != nil {
		// Convert notification to map for Elasticsearch
		notificationMap := map[string]interface{}{}
		notificationBytes, _ := json.Marshal(notification)
		json.Unmarshal(notificationBytes, &notificationMap)

		// Store in Elasticsearch
		_, err = np.ESRepository.Create("device-sphere-notifications", notificationMap)
		if err != nil {
			log.Error("Failed to store notification in Elasticsearch: %v", err)
		} else {
			log.Debug("Notification stored in Elasticsearch: %s", notification.ID)
		}
	}
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
			go np.ruleDataTransformer(string(msg.Data))
			// case msg := <-processedDataTopicChan:
			// 	log.Debug("Received message topic: processed_data: %s", msg.Data)
		}
	}
}
