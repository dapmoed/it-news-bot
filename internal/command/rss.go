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
	textNotFound     = "Пользователь не найден"
)

type SubscribeCallbackData struct {
	RssId uint `json:"rssId"`
}

type RssCommand struct {
	rssRepo          db.RssRepoI
	subscriptionRepo db.SubscriptionRepoI
	userRepo         db.UsersRepoI
	logger           *zap.Logger
	templates        *template.Templates
	bot              *tgbotapi.BotAPI
}

type RssCommandParam struct {
	RssRepo          db.RssRepoI
	SubscriptionRepo db.SubscriptionRepoI
	UserRepo         db.UsersRepoI
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
		userRepo:         param.UserRepo,
	}
}

func (r *RssCommand) List(ctx *chains.Context) {
	errorMessage := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, defaultTextError)
	errorNotFound := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, textNotFound)

	user, err := r.userRepo.GetUser(ctx.Update.Message.From.ID)
	if err != nil {
		r.logger.Error("error get user", zap.Error(err))
		r.bot.Send(errorMessage)
		return
	}

	if user == nil {
		r.bot.Send(errorNotFound)
		return
	}

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

		var numericKeyboardSubscribe = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подписаться", chains.NewCallbackData("subscribe", SubscribeCallbackData{
					RssId: rss.ID,
				}).JSON()),
			),
		)
		var numericKeyboardUnSubscribe = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отписаться", chains.NewCallbackData("unsubscribe", SubscribeCallbackData{
					RssId: rss.ID,
				}).JSON()),
			),
		)

		isSubscribe, err := r.subscriptionRepo.IsSubscribe(rss.ID, user.ID)
		if isSubscribe {
			msg.ReplyMarkup = numericKeyboardUnSubscribe
		} else {
			msg.ReplyMarkup = numericKeyboardSubscribe
		}

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
		user, err := r.userRepo.GetUser(ctx.Update.CallbackQuery.From.ID)
		if err != nil {
			r.logger.Error("error get user", zap.Error(err))
			msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Извините. Произошла ошибка")
			r.bot.Send(msg)
			return
		}

		if user == nil {
			msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Пользователь не найден")
			r.bot.Send(msg)
			return
		}

		err = r.subscriptionRepo.Add(user.ID, data.RssId)
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

func (r *RssCommand) UnSubscribeCallback(ctx *chains.Context, data interface{}) {
	s, err := chains.UnmarshalCallbackData(data, SubscribeCallbackData{})
	if err != nil {
		r.logger.Error("error UnmarshalCallbackData", zap.Error(err))
		msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Извините. Произошла ошибка")
		r.bot.Send(msg)
		return
	}

	if data, ok := s.(*SubscribeCallbackData); ok {
		user, err := r.userRepo.GetUser(ctx.Update.CallbackQuery.From.ID)
		if err != nil {
			r.logger.Error("error get user", zap.Error(err))
			msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Извините. Произошла ошибка")
			r.bot.Send(msg)
			return
		}

		if user == nil {
			msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Пользователь не найден")
			r.bot.Send(msg)
			return
		}

		err = r.subscriptionRepo.Remove(user.ID, data.RssId)
		if err != nil {
			r.logger.Error("error add subscription", zap.Error(err))
			msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Извините. Произошла ошибка")
			r.bot.Send(msg)
			return
		}
	}

	msg := tgbotapi.NewMessage(ctx.Update.CallbackQuery.Message.Chat.ID, "Вы успешно отписаны от ленты")
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

	err = r.rssRepo.Add(u.String(), feed.Title, feed.Description)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Проблема с добавлением URL")
		r.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "URL успешно добавлен")
	r.bot.Send(msg)
	ctx.Chain.End()
}
