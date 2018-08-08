package model

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

func (user *User) Columns() string {
	return "id, name, email, created_at"
}

func (user *User) Fields(u *User) []interface{} {
	return []interface{}{&u.ID, &u.Name, &u.Email, &u.CreatedAt}
}
