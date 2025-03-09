package main

import (
	config "factorysphere/devicesphere/engines/notificationengine/pkg/config"
	notificationprocessor "factorysphere/devicesphere/engines/notificationengine/pkg/processors/notificationprocessor"
	log "factorysphere/devicesphere/engines/notificationengine/pkg/utils/loggers"
)


func main() {
	log.Info("Starting the application")
	log.Info("Strating application...")
	cfg, cfgErr := config.GetConfig()
	if cfgErr != nil {
		log.Error("Unable to read config: %v", cfgErr)
	}

	notificationProcessor := notificationprocessor.NewNotificationProcessor(cfg)
	notificationProcessor.ProcessData()
}
