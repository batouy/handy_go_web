package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	stmt := "insert into users (name, email, hashed_password, created) values(?, ?, ?, ?)"
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, nil
	}
	res, err := m.DB.Exec(stmt, name, email, hashed_pwd, time.Now())
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := "select id, hashed_password from users where email=?"
	row := m.DB.QueryRow(stmt, email)

	var id int
	var pwd []byte

	err := row.Scan(&id, &pwd)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(pwd, []byte(password))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Exists(id int) bool {
	return false
}
