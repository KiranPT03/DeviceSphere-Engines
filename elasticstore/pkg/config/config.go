package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger Logger `yaml:"logger"`
	Kafka  Kafka  `yaml:"kafka"`
}

type Logger struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filepath"`
}

type Kafka struct {
	BootstrapServers string `yaml:"bootstrap_servers"`
	InletTopic       string `yaml:"inlet_topic"`
	ConsumerGroup    string `yaml:"consumer_group"`
}

func GetConfig() (*Config, error) {
	yamlFile, err := os.ReadFile("./pkg/config/application.yaml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &config, nil
}
