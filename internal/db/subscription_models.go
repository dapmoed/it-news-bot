package db

import (
	"gorm.io/gorm"
)

type SubscriptionRepoI interface {
}

type Subscription struct {
	gorm.Model
	Id     int64 `gorm_db:"primaryKey"`
	RssId  int64
	UserId int64
}
