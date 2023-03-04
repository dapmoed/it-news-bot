package db

import (
	"gorm.io/gorm"
	"time"
)

type NotifyRepoI interface {
}

type Notify struct {
	gorm.Model
	Id       int64 `gorm_db:"primaryKey"`
	UserId   int64
	LastTime time.Time
}
