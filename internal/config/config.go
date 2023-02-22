package config

import (
	"github.com/caarlos0/env/v7"
)

type Config struct {
	TelegramBotApi telegramBotApi
	DB             db
}

type telegramBotApi struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN"`
}

type db struct {
	PathDB string `env:"SQLITE_FILE_PATH" envDefault:"./data/bot.db"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
