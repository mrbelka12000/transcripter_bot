package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	PromtPort     string `envconfig:"PROMT_PORT" required:"true"`
	MongoDBURL    string `envconfig:"MONGODB_URL" required:"true"`
	Items         string `envconfig:"ITEMS" default:"defaultCollection"`
	AssemblyKey   string `envconfig:"ASSEMBLY_KEY" required:"false"`
}

func LoadConfig(prefix string) (out Config, err error) {

	if err := godotenv.Load(); err != nil {
		return out, fmt.Errorf("failed to load env: %w", err)
	}

	if err := envconfig.Process(prefix, &out); err != nil {
		return out, fmt.Errorf("failed to process config: %w", err)
	}

	return out, nil
}
