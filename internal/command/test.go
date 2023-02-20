package command

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
)

type Test struct {
	bot       *tgbotapi.BotAPI
	usersRepo db.UsersRepoI
}

func NewTest(bot *tgbotapi.BotAPI, usersRepo db.UsersRepoI) *Test {
	return &Test{
		bot:       bot,
		usersRepo: usersRepo,
	}
}

func (t *Test) Start(ctx *chains.Context) {
	defer ctx.Chain.Next()

	_, strg := ctx.Get("strg")
	if strg == nil {
		ctx.Set("strg", ctx.Update.Message.Text)
	} else {
		ctx.Set("strg", fmt.Sprintf("%s-%s", strg.(string), ctx.Update.Message.Text))
	}

	_, answer := ctx.Get("strg")
	if answer != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, answer.(string))
		t.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "NOT KEY")
	t.bot.Send(msg)

}

func (t *Test) End(ctx *chains.Context) {
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "End of script")
	t.bot.Send(msg)
	ctx.Chain.End()
}
