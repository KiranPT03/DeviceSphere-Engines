package processors

import (
	config "github.com/kiranpt03/factorysphere/devicesphere/engines/elasticstore/pkg/config"
	bitnamikafka "github.com/kiranpt03/factorysphere/devicesphere/engines/elasticstore/pkg/connectors/bitnamikafka"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/elasticstore/pkg/utils/logger"
)

func ProcessData() {
	log.Info("Processing data...")
	configuration, _ := config.GetConfig()
	log.Debug("Configuration filepath: %s", configuration.Logger.FilePath)
	consumer := bitnamikafka.GetKafkaConsumer()
	myTopicChan := consumer.Consume("my_topic")
	testTopicChan := consumer.Consume("test_topic")

	for {
		select {
		case msg := <-myTopicChan:
			log.Debug("Received message: %s", msg)
		case msg := <-testTopicChan:
			log.Debug("Received message: %s", msg)
		}
	}
}
