package db

import "time"

type UsersRepoI interface {
	Close()
	Init() error
	GetUser(id int64) (User, error)
	AddUser(id int64, userName string) error
	UpdateUser(user User) error
}

type User struct {
	id       int64
	UserName string
	LastTime time.Time
}
