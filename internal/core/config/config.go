package core_config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
}

func NewConfig() (*Config, error) {
	timeZone := os.Getenv("TIME_ZONE")
	if timeZone == "" {
		timeZone = "UTC"
	}

	zone, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, fmt.Errorf("load time zone: %s: %w", timeZone, err)
	}

	return &Config{
		TimeZone: zone,
	}, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to create config: %w", err))
	}

	return cfg
}
