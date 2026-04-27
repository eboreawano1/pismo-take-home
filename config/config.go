package config


import (
	"fmt"
	"os"
)

func Load() (Config, error) {
	config := Config { DatabaseURL : os.Getenv("DATABASE_URL"),}

	if config.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is missing")
	}

	return config, nil
}

type Config struct {
	DatabaseURL string
}

