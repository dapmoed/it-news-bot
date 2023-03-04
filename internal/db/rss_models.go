package db

import "gorm.io/gorm"

type RssRepoI interface {
	List() ([]Rss, error)
	Add(url, name string) error
}

type Rss struct {
	gorm.Model
	Id          int64 `gorm_db:"primaryKey"`
	Url         string
	Name        string
	Description string
}
