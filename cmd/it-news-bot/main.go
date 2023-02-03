package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"it-news-bot/internal/config"
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
	bot, err := tgbotapi.NewBotAPI(conf.TelegramBotApi.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	workers := worker.New(bot, updates, 10)
	workers.Init()
	workers.Wait()

	//fp := gofeed.NewParser()
	//feed, _ := fp.ParseURL("https://habr.com/ru/rss/all/all/?fl=ru")
	//for _, v := range feed.Items {
	//	fmt.Println(v.Categories)
	//}
}
