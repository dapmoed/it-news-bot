package notify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"it-news-bot/internal/db"
	"time"
)

type Notifier struct {
	userRepo         db.UsersRepoI
	subscriptionRepo db.SubscriptionRepoI
	bot              *tgbotapi.BotAPI
	logger           *zap.Logger
}

type Param struct {
	UserRepo         db.UsersRepoI
	SubscriptionRepo db.SubscriptionRepoI
	Bot              *tgbotapi.BotAPI
	Logger           *zap.Logger
}

func New(param Param) *Notifier {
	return &Notifier{
		userRepo:         param.UserRepo,
		bot:              param.Bot,
		logger:           param.Logger,
		subscriptionRepo: param.SubscriptionRepo,
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

			_ = users
		}
	}
}
