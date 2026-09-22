package core_verification

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BaseURL  string        `envconfig:"BASE_URL" required:"true"`
	TokenTTL time.Duration `envconfig:"TOKEN_TTL" default:"24h"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("EMAIL_VERIFICATION", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get email verification config: %w", err)
		panic(err)
	}

	return config
}
