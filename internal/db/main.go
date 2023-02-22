package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(fileName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
