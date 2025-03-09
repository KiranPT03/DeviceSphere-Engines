package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   Logger   `yaml:"logger"`
	Nats     Nats     `yaml:"nats"`
	Postgres Postgres `yaml:"postgres"`
}

type Logger struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filepath"`
}

type Nats struct {
	Server        string `yaml:"server"`
	Stream        string `yaml:"stream"`
	InletSubject  string `yaml:"inlet_subject"`
	OutletSubject string `yaml:"outlet_subject"`
	ConsumerGroup string `yaml:"consumer_group"`
	Durable       string `yaml:"durable"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
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
