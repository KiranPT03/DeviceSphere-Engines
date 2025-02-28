package dataprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	config "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/config"
	bitnamikafka "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/connectors/bitnamikafka"
	postgres "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/databases/postgres"
	models "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/models"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/utils/loggers"
)

type DataProcessor struct {
	PGRepository  *postgres.PostgreSQLRepository
	KafkaConsumer *bitnamikafka.KafkaConsumer
	KafkaProducer *bitnamikafka.KafkaProducer
}

func NewDataProcessor(config *config.Config) *DataProcessor {
	// Create a new PostgreSQL repository instance
	repository, err := postgres.NewPostgreSQLRepository(config)
	if err != nil {
		panic(err)
	}

	kafkaConfig := bitnamikafka.KafkaConsumerConfig{
		Brokers:  []string{config.Kafka.BootstrapServers},
		GroupID:  config.Kafka.ConsumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Topics:   []string{config.Kafka.InletTopic, "processed_data"},
	}
	consumer := bitnamikafka.NewKafkaConsumer(kafkaConfig)

	kafkaProducerConfig := bitnamikafka.KafkaProducerConfig{
		Brokers: []string{config.Kafka.BootstrapServers},
		Topic:   config.Kafka.OutletTopic,
	}

	producer := bitnamikafka.NewKafkaProducer(kafkaProducerConfig)

	return &DataProcessor{
		PGRepository:  repository,
		KafkaConsumer: consumer,
		KafkaProducer: producer,
	}
}

func (dp *DataProcessor) convertResultsToDevice(results []map[string]interface{}) (*models.Device, error) {
	if len(results) == 0 {
		return nil, nil
	}
	device := &models.Device{Properties: []models.Property{}}
	propertyMap := make(map[string]models.Property)

	for _, row := range results {
		// Device fields
		if device.ID == "" {
			device.ID = row["deviceid"].(string)
			device.ReferenceID = row["devicereferenceid"].(string)
			device.Type = row["devicetype"].(string)
			device.DeviceName = row["devicename"].(string)

			if createdAt, ok := row["devicecreatedat"].(time.Time); ok {
				device.CreatedAt = createdAt.Format(time.RFC3339) // Format time to string
			} else {
				device.CreatedAt = "" // Handle null time.
			}

			device.State = row["devicestate"].(string)
			device.Location = row["devicelocation"].(string)
			device.Status = row["devicestatus"].(string)
			device.Customer = row["devicecustomer"].(string)
			device.Site = row["devicesite"].(string)
		}

		// Property fields
		property := models.Property{
			ID:          row["propertyid"].(string),
			ReferenceID: row["propertyreferenceid"].(string),
			Name:        row["propertyname"].(string),
			Unit:        row["propertyunit"].(string),
			State:       row["propertystate"].(string),
			Status:      row["propertystatus"].(string),
			DataType:    row["propertydatatype"].(string),
			Value:       row["propertyvalue"].(string),
			Threshold:   row["propertythreshold"].(string),
		}

		propertyMap[property.ID] = property
	}
	for _, property := range propertyMap {
		device.Properties = append(device.Properties, property)
	}

	return device, nil
}

func (dp *DataProcessor) getDeviceData(deviceRefId, propertyRefId string) ([]map[string]interface{}, error) {
	query := `
        SELECT
            d.id AS deviceid,
            d.reference_id AS devicereferenceid,
            d.type AS devicetype,
            d.device_name AS devicename,
            d.created_at AS devicecreatedat,
            d.state AS devicestate,
            d.location AS devicelocation,
            d.status AS devicestatus,
            d.customer AS devicecustomer,
            d.site AS devicesite,
            p.id AS propertyid,
            p.reference_id AS propertyreferenceid,
            p.name AS propertyname,
            p.unit AS propertyunit,
            p.state AS propertystate,
            p.status AS propertystatus,
            p.data_type AS propertydatatype,
            p.value AS propertyvalue,
            p.threshold AS propertythreshold
        FROM
            devices d
        JOIN
            properties p ON d.id = p.device_id
        WHERE
            d.reference_id = $1 AND p.reference_id = $2;
`

	results, err := dp.PGRepository.ExecuteQuery(query, deviceRefId, propertyRefId)
	if err != nil {
		log.Error("Error while executing the query: %v", err)
		return nil, err
	}

	fmt.Printf("%+v\n", results)
	return results, nil
}

func (dp *DataProcessor) dataTransformer(data string) {

	var rawData models.RawData
	err := json.Unmarshal([]byte(data), &rawData)
	if err != nil {
		log.Error("Error unmarshalling JSON: %v", err)
	}

	log.Debug("DeviceRefID: %s\n", rawData.DeviceRefID)
	for _, prop := range rawData.Properties {
		log.Debug("  PropertyRefID: %s, Value: %s\n", prop.PropertyRefID, prop.Value)
		result, devErr := dp.getDeviceData(rawData.DeviceRefID, prop.PropertyRefID)
		if devErr != nil {
			log.Error("Error unmarshalling JSON: %v", devErr)
		}
		log.Debug("Received device Data: %v", result)
		deviceModel, modelEr := dp.convertResultsToDevice(result)
		if modelEr != nil {
			log.Error("Error unmarshalling JSON: %v", modelEr)
		}
		log.Debug("Device model received: %v", deviceModel)

		deviceModelBytes, marshalErr := json.Marshal(deviceModel)
		if marshalErr != nil {
			log.Error("Error marshalling device model: %v", marshalErr)
			continue // Skip to the next property if marshalling fails
		}

		ctx := context.Background() // Or use a context with timeout/cancellation
		producer := dp.KafkaProducer
		defer producer.Close()
		err = producer.Produce(ctx, []byte(deviceModel.ReferenceID), deviceModelBytes)
		if err != nil {
			log.Error("Error producing message to Kafka: %v", err)
		}
	}

}

func (dp *DataProcessor) ProcessData() {
	log.Info("Processing the data")
	appConfig, _ := config.GetConfig()
	log.Debug("Printing config %v", appConfig)
	consumer := dp.KafkaConsumer
	defer consumer.Close()
	rawDataTopicChan := consumer.Consume("raw_data")
	processedDataTopicChan := consumer.Consume("processed_data")

	for {
		select {
		case msg := <-rawDataTopicChan:
			dp.dataTransformer(string(msg))
		case msg := <-processedDataTopicChan:
			log.Debug("Received message topic: processed_data: %s", msg)
		}
	}
}
