package db

import (
	"gorm.io/gorm"
)

type SubscriptionRepoI interface {
	Add(userId, rssId int64) error
}

type Subscription struct {
	gorm.Model
	Id     int64 `gorm_db:"primaryKey"`
	RssId  int64
	UserId int64
}
