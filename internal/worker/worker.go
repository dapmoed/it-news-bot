package worker

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
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
			if update.Message == nil {
				continue
			}

			if !update.Message.IsCommand() {
				continue
			}

			handlerCommand := NewHandleCommand(w.Bot, w.config.UsersRepo, update)
			switch update.Message.Command() {
			case "start":
				handlerCommand.Start()
			case "news":
				handlerCommand.ListNews()
			case "test":
				handlerCommand.Test()
			default:
				handlerCommand.Unknown()
			}
		}
	}
}
func (w *Worker) Wait() {
	w.wgGroup.Wait()
}

type HandleCommand struct {
	UsersRepo *db.Repository
	update    tgbotapi.Update
	bot       *tgbotapi.BotAPI
}

func NewHandleCommand(api *tgbotapi.BotAPI, repository *db.Repository, update tgbotapi.Update) *HandleCommand {
	return &HandleCommand{
		update:    update,
		UsersRepo: repository,
		bot:       api,
	}
}

func (h *HandleCommand) Start() error {
	user, err := h.UsersRepo.GetUser(h.update.Message.From.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			username := fmt.Sprintf("%s.%s", h.update.Message.From.FirstName, h.update.Message.From.LastName)
			if h.update.Message.From.UserName != "" {
				username = fmt.Sprintf("%s( @%s )", username, h.update.Message.From.UserName)
			}
			err := h.UsersRepo.AddUser(h.update.Message.From.ID, username)
			if err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(h.update.Message.Chat.ID, fmt.Sprintf("Будем знакомы, %s", username))
			h.bot.Send(msg)
			return nil
		}
		fmt.Println(err)
		msg := tgbotapi.NewMessage(h.update.Message.Chat.ID, "Ошибка поиска пользователя")
		h.bot.Send(msg)
		return nil
	}

	msg := tgbotapi.NewMessage(h.update.Message.Chat.ID, fmt.Sprintf("Привет, %s. Последний раз мы виделись с тобой %s назад", user.UserName, time.Now().Sub(user.LastTime).String()))
	h.bot.Send(msg)

	err = h.UsersRepo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (h *HandleCommand) Unknown() {
	msg := tgbotapi.NewMessage(h.update.Message.Chat.ID, "Я не знаю такой команды")
	h.bot.Send(msg)
}

func (h *HandleCommand) ListNews() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://habr.com/ru/rss/all/all/?fl=ru")
	for _, v := range feed.Items {
		msg := tgbotapi.NewMessage(h.update.Message.Chat.ID, v.Link)
		h.bot.Send(msg)
	}
}

func (h *HandleCommand) Test() {
	msg := tgbotapi.NewMessage(h.update.Message.Chat.ID, "TEST")
	h.bot.Send(msg)
}
