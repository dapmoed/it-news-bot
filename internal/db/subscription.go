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

func (s *SubscriptionRepository) Add(userId, rssId uint) error {
	tx := s.db.Create(&Subscription{
		RssID:  rssId,
		UserID: userId,
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *SubscriptionRepository) Remove(userID, rssID uint) error {
	tx := s.db.Where(Subscription{RssID: rssID, UserID: userID}).Delete(&Subscription{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (s *SubscriptionRepository) IsSubscribe(rssID uint, userID uint) (bool, error) {
	var subs Subscription
	tx := s.db.Where(&Subscription{RssID: rssID, UserID: userID}).Find(&subs)
	if tx.Error != nil {
		return false, tx.Error
	}
	if tx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
