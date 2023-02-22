package worker

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := w.Bot.Request(callback); err != nil {
					panic(err)
				}

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				if _, err := w.Bot.Send(msg); err != nil {
					panic(err)
				}
				continue
			}
			userId := update.Message.From.ID

			if update.Message.IsCommand() {
				chain, err := w.config.ChainsPool.GetChain(update.Message.Command())
				if err != nil {
					if err == chains.ErrNotFound {
						Unknown(w.Bot, update)
						continue
					}
					// TODO log errors
					continue
				}
				w.config.StorageSession.Add(userId, chain)
			}

			userSession, err := w.config.StorageSession.Get(userId)
			if err != nil {
				if err == sessions.ErrNotFound {
					// TODO log error
					continue
				}
				// TODO log error
				continue
			}
			userSession.Extend()
			userSession.GetChain().Call(update)
		}
	}
}
func (w *Worker) Wait() {
	w.wgGroup.Wait()
}

func Unknown(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю такой команды")
	bot.Send(msg)
}
