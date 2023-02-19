package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
)

type News struct {
	bot       *tgbotapi.BotAPI
	usersRepo db.UsersRepoI
}

func NewNews(bot *tgbotapi.BotAPI, usersRepo db.UsersRepoI) *News {
	return &News{
		bot:       bot,
		usersRepo: usersRepo,
	}
}

func (n *News) Start(ctx chains.Context) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://habr.com/ru/rss/all/all/?fl=ru")
	for _, v := range feed.Items {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, v.Link)
		n.bot.Send(msg)
	}
}
