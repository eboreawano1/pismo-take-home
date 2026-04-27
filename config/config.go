package config


import (
	"fmt"
	"os"
)

func Load() (Config, error) {
	config := Config { 
		DatabaseURL : os.Getenv("DATABASE_URL"),
		KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
	}

	if config.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is missing")
	}

	if config.KafkaBrokers == "" || config.KafkaGroupID == "" || config.KafkaTopic == ""  {
		return Config{}, fmt.Errorf("kafka configuration is incomplete")
	}

	return config, nil
}

type Config struct {
	KafkaTopic   string
	KafkaGroupID string
	KafkaBrokers string
	DatabaseURL string
}

