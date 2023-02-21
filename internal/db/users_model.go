package db

import (
	"gorm.io/gorm"
	"time"
)

type UsersRepoI interface {
	Close()
	Init() error
	GetUser(id int64) (*User, error)
	AddUser(id int64, userName string) error
	UpdateUser(user *User) error
	UpdateLastTime(user *User) error
}

type User struct {
	gorm.Model
	Id       int64 `gorm_db:"primaryKey"`
	UserName string
	LastTime time.Time
}
