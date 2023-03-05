package db

import (
	"gorm.io/gorm"
)

type NotifyRepoI interface {
	Add(userID uint, url string) error
	Get(userID uint, url string) (*Notify, error)
}

type Notify struct {
	gorm.Model
	UserID  uint
	URLHash string
}
