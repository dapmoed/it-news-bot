package db

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (*UsersRepository, error) {
	// Migrate the schema
	err := db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return &UsersRepository{
		db: db,
	}, nil
}

func (r *UsersRepository) GetUser(id int64) (*User, error) {
	var user User
	result := r.db.Find(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	return &user, nil
}

func (r *UsersRepository) AddUser(id int64, userName string) error {
	user := &User{
		Id:       id,
		UserName: userName,
		LastTime: time.Now(),
	}
	r.db.Create(&user)
	return nil
}

func (r *UsersRepository) UpdateUser(user *User) error {
	r.db.Save(&user)
	return nil
}

func (r *UsersRepository) UpdateLastTime(user *User) error {
	user.LastTime = time.Now()
	r.db.Save(&user)
	return nil
}
