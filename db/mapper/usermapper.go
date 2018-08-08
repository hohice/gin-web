package mapper

import (
	"database/sql"
	"fmt"

	"github.com/hohice/gin-web/db/model"
)

type UserMapper struct {
	DB *sql.DB
}

// CreateUser inserts a user record into the database and returns a User object
func (um *UserMapper) CreateUser(name, email, pw string) (*model.User, error) {
	u := &model.User{}
	err := um.DB.QueryRow(fmt.Sprintf(`
		INSERT INTO users(name, email, password)
		VALUES($1, $2, crypt($3, gen_salt('bf', 8)))
		RETURNING %s
	`, u.Columns()), name, email, pw).Scan(u.Fields(u)...)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindUserByEmailPassword finds a user with the given email and password
func (um *UserMapper) FindUserByEmailPassword(email, pw string) (*model.User, error) {
	return um.SelectUser(
		"lower(email) = lower($1) AND password = crypt($2, password)",
		email,
		pw,
	)
}

// FindUser finds a user with the given id or returns nil
func (um *UserMapper) FindUser(id string) (*model.User, error) {
	return um.SelectUser("id = $1", id)
}

func (um *UserMapper) SelectUser(c string, v ...interface{}) (*model.User, error) {
	u := &model.User{}
	err := um.DB.QueryRow(
		fmt.Sprintf(`SELECT %s FROM users WHERE %s`, u.Columns(), c),
		v...,
	).Scan(u.Fields(u)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}
