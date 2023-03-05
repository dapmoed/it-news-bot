package db

import (
	"gorm.io/gorm"
	"time"
)

type UsersRepoI interface {
	GetUser(id int64) (*User, error)
	AddUser(tgUserID int64, tgChatID int64, userName string) error
	UpdateUser(user *User) error
	UpdateLastTime(user *User) error
	List() ([]User, error)
}

type User struct {
	gorm.Model
	TgID          int64
	TgChatID      int64
	Name          string
	LastTime      time.Time
	Subscriptions []Subscription
}
