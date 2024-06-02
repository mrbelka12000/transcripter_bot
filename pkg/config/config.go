package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken  string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	PromtPort      string `envconfig:"PROMT_PORT" required:"true"`
	AssemblyKey    string `envconfig:"ASSEMBLY_KEY" required:"true"`
	MongoDBURL     string `envconfig:"MONGODB_URL" required:"true"`
	DBName         string `envconfig:"DB_NAME" default:"go-mongo"`
	CollectionName string `envconfig:"COLLECTION_NAME" default:"defaultCollection"`
}

func LoadConfig(prefix string) (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return cfg, nil
}
