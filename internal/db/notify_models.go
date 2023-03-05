package db

import (
	"gorm.io/gorm"
	"time"
)

type NotifyRepoI interface {
	Add(userID uint) error
	Get(userID uint) (*Notify, error)
	Update(userID uint) error
}

type Notify struct {
	gorm.Model
	UserID   uint
	LastTime time.Time
}
