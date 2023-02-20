package command

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
	"time"
)

type Start struct {
	bot       *tgbotapi.BotAPI
	usersRepo db.UsersRepoI
}

func NewStart(bot *tgbotapi.BotAPI, usersRepo db.UsersRepoI) *Start {
	return &Start{
		bot:       bot,
		usersRepo: usersRepo,
	}
}

func (c *Start) Start(ctx *chains.Context) {
	user, err := c.usersRepo.GetUser(ctx.Update.Message.From.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			username := fmt.Sprintf("%s.%s", ctx.Update.Message.From.FirstName, ctx.Update.Message.From.LastName)
			if ctx.Update.Message.From.UserName != "" {
				username = fmt.Sprintf("%s( @%s )", username, ctx.Update.Message.From.UserName)
			}
			err := c.usersRepo.AddUser(ctx.Update.Message.From.ID, username)
			if err != nil {
				// TODO log
				return
			}
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, fmt.Sprintf("Будем знакомы, %s", username))
			c.bot.Send(msg)
			return
		}
		fmt.Println(err)
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Ошибка поиска пользователя")
		c.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, fmt.Sprintf("Привет, %s. Последний раз мы виделись с тобой %s назад", user.UserName, time.Now().Sub(user.LastTime).String()))
	c.bot.Send(msg)

	err = c.usersRepo.UpdateUser(user)
	if err != nil {
		// TODO logs
		return
	}
	return
	//ctx.Chain.Next()
}
