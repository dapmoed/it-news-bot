package command

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
	"os"
	"text/template"
)

const (
	defaultTextError = "Извините у нас проблемы"
)

type RssCommand struct {
	bot     *tgbotapi.BotAPI
	rssRepo db.RssRepoI
	logger  *zap.Logger
}

func NewCommandRss(bot *tgbotapi.BotAPI, rssRepo db.RssRepoI, logger *zap.Logger) *RssCommand {
	return &RssCommand{
		bot:     bot,
		rssRepo: rssRepo,
		logger:  logger,
	}
}

func (r *RssCommand) List(ctx *chains.Context) {
	errorMessage := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, defaultTextError)

	rss, err := r.rssRepo.List()
	if err != nil {
		r.logger.Error("error list rss", zap.Error(err))
		r.bot.Send(errorMessage)
		return
	}

	fmt.Println(os.Getwd())

	tmpl, err := template.ParseFiles("data/list_rss.tmpl")
	if err != nil {
		r.logger.Error("error parse template", zap.Error(err))
		r.bot.Send(errorMessage)
		return
	}

	var textMessage bytes.Buffer

	err = tmpl.Execute(&textMessage, rss)
	if err != nil {
		r.logger.Error("error execute template", zap.Error(err))
		r.bot.Send(errorMessage)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, textMessage.String())

	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Add", "Add"),
		),
	)
	msg.ReplyMarkup = numericKeyboard
	_, err = r.bot.Send(msg)
	if err != nil {
		r.logger.Error("error send message", zap.Error(err))
	}
	ctx.Chain.Next()
}

func (r *RssCommand) Add(ctx *chains.Context) {
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Update.Message.Text)

	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := r.bot.Send(msg)
	if err != nil {
		r.logger.Error("error send message", zap.Error(err))
	}
}
