package core_verification

import (
	"errors"
	"fmt"
	"net/url"
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

	if err := config.validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return config, nil
}

func (c Config) validate() error {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("EMAIL_VERIFICATION_BASE_URL is not a valid URL: %w", err)
	}

	if (base.Scheme != "http" && base.Scheme != "https") || base.Host == "" {
		return errors.New(
			"EMAIL_VERIFICATION_BASE_URL must be an absolute http(s) URL, " +
				"e.g. https://api.example.com/api/v1/auth/verify-email",
		)
	}

	if c.TokenTTL <= 0 {
		return errors.New("EMAIL_VERIFICATION_TOKEN_TTL must be positive")
	}

	return nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get email verification config: %w", err)
		panic(err)
	}

	return config
}
