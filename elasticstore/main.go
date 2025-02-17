package main

import (
	processor "github.com/kiranpt03/factorysphere/devicesphere/engines/elasticstore/pkg/processors"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/elasticstore/pkg/utils/logger"
)

func main() {
	log.Info("Starting elastic store engine")
	processor.ProcessData()
}
