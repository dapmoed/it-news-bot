package gorm_db

import (
	"database/sql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	db2 "it-news-bot/internal/db"
	"time"
)

type UsersRepository struct {
	db *gorm.DB
}

func New(fileName string) (*UsersRepository, error) {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&db2.User{})
	if err != nil {
		return nil, err
	}

	return &UsersRepository{
		db: db,
	}, nil
}

func (r *UsersRepository) Close() {
	db, err := r.db.DB()
	if err != nil {
		// TODO LOG
	}
	if err := db.Close(); err != nil {
		//TODO LOG
	}
}

func (r *UsersRepository) Init() error {
	return nil
}

func (r *UsersRepository) GetUser(id int64) (*db2.User, error) {
	var user db2.User
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
	user := &db2.User{
		Id:       id,
		UserName: userName,
		LastTime: time.Now(),
	}
	r.db.Create(&user)
	return nil
}

func (r *UsersRepository) UpdateUser(user *db2.User) error {
	r.db.Save(&user)
	return nil
}

func (r *UsersRepository) UpdateLastTime(user *db2.User) error {
	user.LastTime = time.Now()
	r.db.Save(&user)
	return nil
}
