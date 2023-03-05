package db

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

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

func (n *NotifyRepository) Add(userID uint) error {
	notify := &Notify{
		LastTime: time.Now(),
		UserID:   userID,
	}
	n.db.Create(notify)
	return nil
}

func (n *NotifyRepository) Get(userID uint) (*Notify, error) {
	var notify Notify
	result := n.db.First(&notify, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, sql.ErrNoRows
		}
		return nil, result.Error
	}
	return &notify, nil
}

func (n *NotifyRepository) Update(userID uint) error {
	notify, err := n.Get(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			n.Add(userID)
			return nil
		}
		return err
	}

	notify.LastTime = time.Now()
	n.db.Save(notify)
	return nil
}
