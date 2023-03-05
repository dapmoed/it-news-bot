package db

import "gorm.io/gorm"

type RssRepoI interface {
	List() ([]Rss, error)
	Add(url, name, description string) error
	Get(id uint) (*Rss, error)
}

type Rss struct {
	gorm.Model
	Url         string
	Name        string
	Description string
}
