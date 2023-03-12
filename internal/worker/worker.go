package worker

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"it-news-bot/internal/chains"
	"it-news-bot/internal/db"
	"it-news-bot/internal/sessions"
	"sync"
)

type Worker struct {
	wgGroup *sync.WaitGroup
	Count   int
	Chanel  tgbotapi.UpdatesChannel
	Bot     *tgbotapi.BotAPI
	config  Config
}

type Config struct {
	UsersRepo      db.UsersRepoI
	ChainsPool     *chains.Pool
	StorageSession *sessions.StorageSession
	Logger         *zap.Logger
}

func New(botApi *tgbotapi.BotAPI, chanel tgbotapi.UpdatesChannel, config Config, count int) *Worker {
	return &Worker{
		Chanel:  chanel,
		Count:   count,
		Bot:     botApi,
		wgGroup: &sync.WaitGroup{},
		config:  config,
	}
}

func (w *Worker) Init() {
	for i := 0; i < w.Count; i++ {
		w.wgGroup.Add(1)
		go w.Handle(i)
	}
}

func (w *Worker) Handle(name int) {
	defer w.wgGroup.Done()
	for {
		select {
		case update := <-w.Chanel:
			if update.Message == nil {
				//callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				//if _, err := w.Bot.Request(callback); err != nil {
				//	SendError(w.Bot, update.CallbackQuery.Message.Chat.ID)
				//	w.config.Logger.Error("error request Callback", zap.Error(err))
				//	continue
				//}

				userSession, err := w.config.StorageSession.Get(update.CallbackQuery.From.ID)
				if err != nil {
					w.config.Logger.Error("error get user session", zap.Error(err))
					continue
				}

				userSession.Extend()
				err = userSession.GetChain().CallCallback(update)
				if err != nil {
					SendError(w.Bot, update.CallbackQuery.Message.Chat.ID)
					w.config.Logger.Error("error call func callback", zap.Error(err))
					continue
				}
				continue
			}

			userId := update.Message.From.ID

			if update.Message.IsCommand() {
				chain, err := w.config.ChainsPool.GetChain(update.Message.Command())
				if err != nil {
					if err == chains.ErrNotFound {
						SendUnknown(w.Bot, update.Message.Chat.ID)
						w.config.Logger.Error("not found chain", zap.Error(err))
						continue
					}
					w.config.Logger.Error("error get chain", zap.Error(err))
					continue
				}
				w.config.StorageSession.Add(userId, chain)
			}

			userSession, err := w.config.StorageSession.Get(userId)
			if err != nil {
				if err == sessions.ErrNotFound {
					w.config.Logger.Error("not found user session", zap.Error(err))
					continue
				}
				w.config.Logger.Error("error get user session", zap.Error(err))
				continue
			}
			userSession.Extend()
			err = userSession.GetChain().Call(update)
			if err != nil {
				w.config.Logger.Error("error call step chain", zap.Error(err))
				continue
			}
		}
	}
}
func (w *Worker) Wait() {
	w.wgGroup.Wait()
}

func SendUnknown(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Я не знаю такой команды")
	bot.Send(msg)
}

func SendError(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Произошла внутренняя ошибка")
	bot.Send(msg)
}
