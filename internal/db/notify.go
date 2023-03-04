package db

import "gorm.io/gorm"

type NotifyRepository struct {
	db *gorm.DB
}

func NewNotifyRepo(db *gorm.DB) (*NotifyRepository, error) {
	// Migrate the schema
	err := db.AutoMigrate(&Notify{})
	if err != nil {
		return nil, err
	}

	return &NotifyRepository{
		db: db,
	}, nil
}
