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
	HTTPPort       string `envconfig:"HTTP_PORT" default:"5551"`
	BotName        string `envconfig:"BOT_NAME" default:"@chat_transcripter_bot"`
}

func LoadConfig(prefix string) (out Config, err error) {

	godotenv.Load()

	err = envconfig.Process(prefix, &out)
	if err != nil {
		return out, fmt.Errorf("failed to process config: %w", err)
	}

	return out, nil
}
