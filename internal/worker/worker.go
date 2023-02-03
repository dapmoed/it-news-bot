package worker

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type Worker struct {
	wgGroup *sync.WaitGroup
	Count   int
	Chanel  tgbotapi.UpdatesChannel
	Bot     *tgbotapi.BotAPI
}

func New(botApi *tgbotapi.BotAPI, c tgbotapi.UpdatesChannel, count int) *Worker {
	return &Worker{
		Chanel:  c,
		Count:   count,
		Bot:     botApi,
		wgGroup: &sync.WaitGroup{},
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
			fmt.Println("Name bot:", name)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			message, err := w.Bot.Send(msg)
			if err != nil {
				// log
			}
			_ = message
		}
	}
}
func (w *Worker) Wait() {
	w.wgGroup.Wait()
}
