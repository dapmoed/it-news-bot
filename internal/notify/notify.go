package notify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
	"it-news-bot/internal/db"
	"time"
)

type Notifier struct {
	userRepo         db.UsersRepoI
	subscriptionRepo db.SubscriptionRepoI
	rssRepo          db.RssRepoI
	notifyRepo       db.NotifyRepoI
	bot              *tgbotapi.BotAPI
	logger           *zap.Logger
}

type Param struct {
	UserRepo         db.UsersRepoI
	SubscriptionRepo db.SubscriptionRepoI
	NotifyRepo       db.NotifyRepoI
	RssRepo          db.RssRepoI
	Bot              *tgbotapi.BotAPI
	Logger           *zap.Logger
}

func New(param Param) *Notifier {
	return &Notifier{
		userRepo:         param.UserRepo,
		rssRepo:          param.RssRepo,
		bot:              param.Bot,
		logger:           param.Logger,
		subscriptionRepo: param.SubscriptionRepo,
		notifyRepo:       param.NotifyRepo,
	}
}

func (n *Notifier) Run() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			users, err := n.userRepo.List()
			if err != nil {
				n.logger.Error("error get users", zap.Error(err))
				continue
			}

			for _, user := range users {
				for _, subscription := range user.Subscriptions {
					n.Notify(&user, subscription.RssID)
				}
			}
		}
	}
}

func (n *Notifier) Notify(user *db.User, rssId uint) {
	rss, err := n.rssRepo.Get(rssId)
	if err != nil {
		n.logger.Error("error get rss", zap.Error(err))
		return
	}

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(rss.Url)
	for _, v := range feed.Items {
		notify, err := n.notifyRepo.Get(user.ID, v.Link)
		if err != nil {
			n.logger.Error("error get notify", zap.Error(err))
			continue
		}

		if notify == nil {
			msg := tgbotapi.NewMessage(user.TgChatID, v.Link)
			_, err := n.bot.Send(msg)
			if err != nil {
				n.logger.Error("error send message", zap.Error(err))
			}
			n.notifyRepo.Add(user.ID, v.Link)
		}
	}
}
