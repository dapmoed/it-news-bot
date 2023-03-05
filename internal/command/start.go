package command

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
	"time"
)

type StartCommand struct {
	bot       *tgbotapi.BotAPI
	usersRepo db.UsersRepoI
}

func NewCommandStart(bot *tgbotapi.BotAPI, usersRepo db.UsersRepoI) *StartCommand {
	return &StartCommand{
		bot:       bot,
		usersRepo: usersRepo,
	}
}

func (c *StartCommand) Start(ctx *chains.Context) {
	user, err := c.usersRepo.GetUser(ctx.Update.Message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Ошибка поиска пользователя")
		c.bot.Send(msg)
		return
	}

	if user == nil {
		username := fmt.Sprintf("%s.%s", ctx.Update.Message.From.FirstName, ctx.Update.Message.From.LastName)
		if ctx.Update.Message.From.UserName != "" {
			username = fmt.Sprintf("%s( @%s )", username, ctx.Update.Message.From.UserName)
		}
		err := c.usersRepo.AddUser(ctx.Update.Message.From.ID, ctx.Update.Message.Chat.ID, username)
		if err != nil {
			// TODO log
			return
		}
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, fmt.Sprintf("Будем знакомы, %s", username))
		c.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, fmt.Sprintf("Привет, %s. Последний раз мы виделись с тобой %s назад", user.Name, time.Now().Sub(user.LastTime).String()))
	c.bot.Send(msg)

	err = c.usersRepo.UpdateLastTime(user)
	if err != nil {
		// TODO logs
		return
	}
	return
	//ctx.Chain.Next()
}
