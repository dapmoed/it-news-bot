package db

import (
	"gorm.io/gorm"
)

type SubscriptionRepoI interface {
	Add(userID, rssID uint) error
	IsSubscribe(rssID uint, userID uint) (bool, error)
	Remove(userID, rssID uint) error
}

type Subscription struct {
	gorm.Model
	RssID  uint
	UserID uint
}
