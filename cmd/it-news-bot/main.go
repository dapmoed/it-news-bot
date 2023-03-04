package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/command"
	"it-news-bot/internal/config"
	"it-news-bot/internal/db"
	"it-news-bot/internal/sessions"
	"it-news-bot/internal/template"
	"it-news-bot/internal/worker"
	"log"
	"time"
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
	fmt.Println(conf)

	templates, err := template.New(conf.TPL.Path)
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sqlLiteDB, err := db.NewDB(conf.DB.Path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if db, err := sqlLiteDB.DB(); err == nil {
			err := db.Close()
			if err != nil {
				// TODO LOG
			}
		}
	}()

	usersRepo, err := db.NewUserRepo(sqlLiteDB)
	if err != nil {
		panic(err)
	}
	rssRepo, err := db.NewRssRepo(sqlLiteDB)
	if err != nil {
		panic(err)
	}
	subscriptionRepo, err := db.NewSubscriptionRepo(sqlLiteDB)
	if err != nil {
		panic(err)
	}
	_ = subscriptionRepo

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
	startCommand := command.NewCommandStart(bot, usersRepo)
	chainsPool.Command("start",
		chains.NewChain().
			Register(startCommand.Start).SetDurationSession(time.Second),
	)
	newsCommand := command.NewCommandNews(bot, usersRepo)
	chainsPool.Command("news",
		chains.NewChain().
			Register(newsCommand.Start))

	testCommand := command.NewCommandTest(bot, usersRepo)
	chainsPool.Command("test",
		chains.NewChain().
			Register(testCommand.Start).Register(testCommand.End))

	rssCommand := command.NewCommandRss(bot, rssRepo, templates, logger)
	chainsPool.Command("rss",
		chains.NewChain().
			Register(rssCommand.List).RegisterCallback("Add", rssCommand.AddRssCallback))

	chainsPool.Command("rss_add",
		chains.NewChain().
			Register(rssCommand.AddRssUriStepOne).Register(rssCommand.AddRssUriStepTwo))

	workers := worker.New(bot, updates, worker.Config{
		UsersRepo:      usersRepo,
		ChainsPool:     chainsPool,
		StorageSession: sessions.New(),
		Logger:         logger,
	}, 10)
	workers.Init()
	workers.Wait()

}
