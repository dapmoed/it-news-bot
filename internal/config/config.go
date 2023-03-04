package config

import (
	"github.com/caarlos0/env/v7"
)

type Config struct {
	TelegramBotApi telegramBotApi
	DB             db
	TPL            templates
}

type telegramBotApi struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN"`
}

type db struct {
	Path string `env:"SQLITE_FILE_PATH" envDefault:"./data/bot.db"`
}

type templates struct {
	Path string `env:"TEMPLATE_PATH" envDefault:"./template"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
