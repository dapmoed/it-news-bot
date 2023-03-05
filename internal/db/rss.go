package db

import (
	"database/sql"
	"gorm.io/gorm"
)

type RssRepository struct {
	db *gorm.DB
}

func NewRssRepo(db *gorm.DB) (*RssRepository, error) {
	// Migrate the schema
	err := db.AutoMigrate(&Rss{})
	if err != nil {
		return nil, err
	}

	return &RssRepository{
		db: db,
	}, nil
}

func (r *RssRepository) List() ([]Rss, error) {
	rssItems := make([]Rss, 0)
	result := r.db.Find(&rssItems)
	if result.Error != nil {
		return rssItems, result.Error
	}
	return rssItems, nil
}

func (r *RssRepository) Add(url, name, description string) error {
	rss := &Rss{
		Url:         url,
		Name:        name,
		Description: description,
	}
	tx := r.db.Create(&rss)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RssRepository) Get(id uint) (*Rss, error) {
	var rss Rss
	tx := r.db.First(&rss, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, sql.ErrNoRows
		}
		return nil, tx.Error
	}
	return &rss, nil
}
