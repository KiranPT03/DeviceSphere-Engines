package main

import (
	config "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/config"
	ruleprocessor "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/processors/ruleprocessor"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/utils/loggers"
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
