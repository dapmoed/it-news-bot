package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
)

type NewsCommand struct {
	bot       *tgbotapi.BotAPI
	usersRepo db.UsersRepoI
}

func NewCommandNews(bot *tgbotapi.BotAPI, usersRepo db.UsersRepoI) *NewsCommand {
	return &NewsCommand{
		bot:       bot,
		usersRepo: usersRepo,
	}
}

func (n *NewsCommand) Start(ctx *chains.Context) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://habr.com/ru/rss/all/all/?fl=ru")
	for _, v := range feed.Items {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, v.Link)
		n.bot.Send(msg)
	}
}
