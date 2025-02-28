package main

import (
	config "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/config"
	dataprocessor "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/processors/dataprocessor"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/kafkaingest/pkg/utils/loggers"
)

func main()  {
	log.Info("Starting the application.....")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error while loading config, reason: %s", err.Error())
	}
	dataProcessor := dataprocessor.NewDataProcessor(cfg)
	dataProcessor.ProcessData()
}