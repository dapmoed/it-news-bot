package db

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/gorm"
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

func (n *NotifyRepository) Add(userID uint, url string) error {
	notify := &Notify{
		URLHash: GetMd5(url),
		UserID:  userID,
	}
	n.db.Create(notify)
	return nil
}

func (n *NotifyRepository) Get(userID uint, url string) (*Notify, error) {
	var notify Notify
	tx := n.db.Where(&Notify{URLHash: GetMd5(url), UserID: userID}).Find(&notify)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return &notify, nil
}

//func (n *NotifyRepository) Update(userID uint) error {
//	notify, err := n.Get(userID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			n.Add(userID)
//			return nil
//		}
//		return err
//	}
//
//	notify.LastTime = time.Now()
//	n.db.Save(notify)
//	return nil
//}

func GetMd5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
