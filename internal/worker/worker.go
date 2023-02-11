package worker

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"it-news-bot/internal/db"
	"sync"
	"time"
)

type Worker struct {
	wgGroup *sync.WaitGroup
	Count   int
	Chanel  tgbotapi.UpdatesChannel
	Bot     *tgbotapi.BotAPI
	config  Config
}

type Config struct {
	UsersRepo *db.Repository
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
			user, err := w.config.UsersRepo.GetUser(update.Message.From.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					err := w.config.UsersRepo.AddUser(update.Message.From.ID, update.Message.From.UserName)
					if err != nil {
						fmt.Println(err)
						continue
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Будем знакомы, %s", update.Message.From.UserName))
					w.Bot.Send(msg)
					continue
				}
				fmt.Println(err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска пользователя")
				w.Bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s. Последний раз мы виделись с тобой %s назад", user.UserName, time.Now().Sub(user.LastTime).String()))
			w.Bot.Send(msg)

			err = w.config.UsersRepo.UpdateUser(user)
			if err != nil {
				fmt.Println(err)
			}

			//
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID
			//message, err := w.Bot.Send(msg)
			//if err != nil {
			//	// log
			//}
			//_ = message
		}
	}
}
func (w *Worker) Wait() {
	w.wgGroup.Wait()
}
