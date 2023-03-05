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

func (r *UsersRepository) GetUser(tgUserID int64) (*User, error) {
	var user User
	result := r.db.Model(User{TgID: tgUserID}).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, sql.ErrNoRows
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UsersRepository) AddUser(tgUserID int64, tgChatID int64, userName string) error {
	user := &User{
		TgID:     tgUserID,
		TgChatID: tgChatID,
		Name:     userName,
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

func (r *UsersRepository) List() ([]User, error) {
	users := make([]User, 0)
	err := r.db.Model(&User{}).Preload("Subscriptions").Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, nil
}
