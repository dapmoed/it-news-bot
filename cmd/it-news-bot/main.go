package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/command"
	"it-news-bot/internal/config"
	"it-news-bot/internal/db"
	"it-news-bot/internal/sessions"
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
	usersRepo, err := db.New("./data/bot_users.db")
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

	chainsPool := chains.NewPool()
	startCommand := command.NewStart(bot, usersRepo)
	chainsPool.Command("start",
		chains.NewChain().
			Register(startCommand.Start),
	)
	newsCommand := command.NewNews(bot, usersRepo)
	chainsPool.Command("news",
		chains.NewChain().
			Register(newsCommand.Start))

	testCommand := command.NewTest(bot, usersRepo)
	chainsPool.Command("test",
		chains.NewChain().
			Register(testCommand.Start).Register(testCommand.End))

	workers := worker.New(bot, updates, worker.Config{
		UsersRepo:      usersRepo,
		ChainsPool:     chainsPool,
		StorageSession: sessions.New(),
	}, 10)
	workers.Init()
	workers.Wait()

}
