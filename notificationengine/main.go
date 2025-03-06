package main

import (
	config "devicesphere/engines/notificationengine/pkg/config"
	ruleprocessor "devicesphere/engines/notificationengine/pkg/processors/ruleprocessor"
	log "devicesphere/engines/notificationengine/pkg/utils/loggers"
)

func main() {
	log.Info("Strating application...")
	cfg, cfgErr := config.GetConfig()
	if cfgErr != nil {
		log.Error("Unable to read config: %v", cfgErr)
	}

	ruleProcessor := ruleprocessor.NewRuleProcessor(cfg)
	ruleProcessor.ProcessData()
}
