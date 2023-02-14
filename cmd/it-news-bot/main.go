package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"it-news-bot/internal/config"
	"it-news-bot/internal/db"
	"it-news-bot/internal/worker"
	"log"
)

var (
	conf config.Config
)

func init() {
	var err error
	conf, err = config.GetConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	usersRepo, err := db.New("data/bot_users.db")
	if err != nil {
		panic(err)
	}
	defer usersRepo.Close()
	err = usersRepo.Init()
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.TelegramBotApi.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	workers := worker.New(bot, updates, worker.Config{
		UsersRepo: usersRepo,
	}, 10)
	workers.Init()
	workers.Wait()

}
