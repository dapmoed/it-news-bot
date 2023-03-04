package command

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
	"it-news-bot/internal/template"
	"net/url"
)

const (
	defaultTextError = "Извините у нас проблемы"
)

type RssCommand struct {
	bot       *tgbotapi.BotAPI
	rssRepo   db.RssRepoI
	logger    *zap.Logger
	templates *template.Templates
}

func NewCommandRss(bot *tgbotapi.BotAPI, rssRepo db.RssRepoI, templates *template.Templates, logger *zap.Logger) *RssCommand {
	return &RssCommand{
		bot:       bot,
		rssRepo:   rssRepo,
		logger:    logger,
		templates: templates,
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

	var textMessage bytes.Buffer

	err = r.templates.Execute("list_rss", &textMessage, rss)
	if err != nil {
		r.logger.Error("error execute template", zap.Error(err))
		r.bot.Send(errorMessage)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, textMessage.String())

	//var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	//	tgbotapi.NewInlineKeyboardRow(
	//		tgbotapi.NewInlineKeyboardButtonData("Add", "Add"),
	//	),
	//)
	//msg.ReplyMarkup = numericKeyboard
	_, err = r.bot.Send(msg)
	if err != nil {
		r.logger.Error("error send message", zap.Error(err))
	}
}

//func (r *RssCommand) Add(ctx *chains.Context) {
//	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Update.Message.Text)
//
//	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
//	_, err := r.bot.Send(msg)
//	if err != nil {
//		r.logger.Error("error send message", zap.Error(err))
//	}
//}

func (r *RssCommand) AddRssCallback(ctx *chains.Context) {
	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Команда сработала")
	r.bot.Send(msg)
}

func (r *RssCommand) AddRssUriStepOne(ctx *chains.Context) {
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Отправьте в ответном сообщении URI RSS ленты")
	r.bot.Send(msg)
	ctx.Chain.Next()
}

func (r *RssCommand) AddRssUriStepTwo(ctx *chains.Context) {
	u, err := url.ParseRequestURI(ctx.Update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Это не ссылка")
		r.bot.Send(msg)
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u.String())
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Проблемы с получением ленты")
		r.bot.Send(msg)
		return
	}

	err = r.rssRepo.Add(u.String(), feed.Title)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Проблема с добавлением URL")
		r.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "URL успешно добавлен")
	r.bot.Send(msg)
	ctx.Chain.End()
}
