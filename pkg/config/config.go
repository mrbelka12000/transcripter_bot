package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	PromtPort     string `envconfig:"PROMT_PORT" required:"true"`
	AssemblyKey   string `envconfig:"ASSEMBLY_KEY" required:"true"`
	TableName     string `envconfig:"TABLE_NAME" required:"true"`
	HTTPPort      string `envconfig:"HTTP_PORT" default:"5551"`
	BotName       string `envconfig:"BOT_NAME" default:"@chat_transcripter_bot"`
	PGURL         string `envconfig:"PG_URL"  required:"true"`
}

func LoadConfig(prefix string) (out Config, err error) {

	godotenv.Load()

	err = envconfig.Process(prefix, &out)
	if err != nil {
		return out, fmt.Errorf("failed to process config: %w", err)
	}

	return out, nil
}
