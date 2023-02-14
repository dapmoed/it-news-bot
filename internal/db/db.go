package db

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(fileName string) (*Repository, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) Init() error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS users(  id INTEGER PRIMARY KEY AUTOINCREMENT,  username TEXT,  last_time INTEGER);")
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	id       int64
	UserName string
	LastTime time.Time
}

func (r *Repository) GetUser(id int64) (User, error) {
	user := User{}
	rows := r.db.QueryRow("select * from users")
	iTime := int64(0)
	err := rows.Scan(&user.id, &user.UserName, &iTime)
	if err != nil {
		return user, err
	}
	user.LastTime = time.Unix(iTime, 0)
	return user, nil
}

func (r *Repository) AddUser(id int64, userName string) error {
	_, err := r.db.Exec("insert into users (id, username, last_time) values ($1, $2, $3)",
		id, userName, time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateUser(user User) error {
	_, err := r.db.Exec("update users set username = $1, last_time = $2 where id = $3", user.UserName, time.Now().Unix(), user.id)
	if err != nil {
		return err
	}
	return nil
}
