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

type SubscribeCallbackData struct {
	RssId int64 `json:"rssId"`
}

type RssCommand struct {
	rssRepo          db.RssRepoI
	subscriptionRepo db.SubscriptionRepoI
	logger           *zap.Logger
	templates        *template.Templates
	bot              *tgbotapi.BotAPI
}

type RssCommandParam struct {
	RssRepo          db.RssRepoI
	SubscriptionRepo db.SubscriptionRepoI
	Logger           *zap.Logger
	Templates        *template.Templates
	Bot              *tgbotapi.BotAPI
}

func NewCommandRss(param RssCommandParam) *RssCommand {
	return &RssCommand{
		bot:              param.Bot,
		rssRepo:          param.RssRepo,
		logger:           param.Logger,
		templates:        param.Templates,
		subscriptionRepo: param.SubscriptionRepo,
	}
}

func (r *RssCommand) List(ctx *chains.Context) {
	errorMessage := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, defaultTextError)

	rssList, err := r.rssRepo.List()
	if err != nil {
		r.logger.Error("error list rss", zap.Error(err))
		r.bot.Send(errorMessage)
		return
	}

	for _, rss := range rssList {
		var textMessage bytes.Buffer
		err = r.templates.Execute("rss", &textMessage, rss)
		if err != nil {
			r.logger.Error("error execute template", zap.Error(err))
			r.bot.Send(errorMessage)
			return
		}
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, textMessage.String())
		var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подписаться", chains.NewCallbackData("subscribe", SubscribeCallbackData{
					RssId: rss.Id,
				}).JSON()),
			),
		)
		msg.ReplyMarkup = numericKeyboard
		_, err = r.bot.Send(msg)
		if err != nil {
			r.logger.Error("error send message", zap.Error(err))
		}
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

func (r *RssCommand) SubscribeCallback(ctx *chains.Context, data interface{}) {
	s, err := chains.UnmarshalCallbackData(data, SubscribeCallbackData{})
	if err != nil {
		r.logger.Error("error UnmarshalCallbackData", zap.Error(err))
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Извините. Произошла ошибка")
		r.bot.Send(msg)
		return
	}

	if data, ok := s.(*SubscribeCallbackData); ok {
		err := r.subscriptionRepo.Add(ctx.Update.CallbackQuery.From.ID, data.RssId)
		if err != nil {
			r.logger.Error("error add subscription", zap.Error(err))
			msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Извините. Произошла ошибка")
			r.bot.Send(msg)
			return
		}
	}

	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Вы успешно подписаны на ленту")
	r.bot.Send(msg)
	return
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
