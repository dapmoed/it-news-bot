package db

import (
	"database/sql"
	"time"
)

type UsersRepository struct {
	db *sql.DB
}

func New(fileName string) (*UsersRepository, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	return &UsersRepository{
		db: db,
	}, nil
}

func (r *UsersRepository) Close() {
	r.db.Close()
}

func (r *UsersRepository) Init() error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS users(  id INTEGER PRIMARY KEY AUTOINCREMENT,  username TEXT,  last_time INTEGER);")
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) GetUser(id int64) (User, error) {
	user := User{}
	rows := r.db.QueryRow("select * from users where id = $1",
		id)
	iTime := int64(0)
	err := rows.Scan(&user.id, &user.UserName, &iTime)
	if err != nil {
		return user, err
	}
	user.LastTime = time.Unix(iTime, 0)
	return user, nil
}

func (r *UsersRepository) AddUser(id int64, userName string) error {
	_, err := r.db.Exec("insert into users (id, username, last_time) values ($1, $2, $3)",
		id, userName, time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) UpdateUser(user User) error {
	_, err := r.db.Exec("update users set username = $1, last_time = $2 where id = $3",
		user.UserName, time.Now().Unix(), user.id)
	if err != nil {
		return err
	}
	return nil
}
