package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN"`
	PromtPort     string `envconfig:"PROMT_PORT"`
}

func LoadConfig(prefix string) (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return cfg, nil
}
