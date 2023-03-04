package db

import (
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) (*SubscriptionRepository, error) {
	// Migrate the schema
	err := db.AutoMigrate(&Subscription{})
	if err != nil {
		return nil, err
	}

	return &SubscriptionRepository{
		db: db,
	}, nil
}

func (s *SubscriptionRepository) Add(userId, rssId int64) error {
	tx := s.db.Create(&Subscription{
		RssId:  rssId,
		UserId: userId,
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
