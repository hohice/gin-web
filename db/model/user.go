package model

import "time"

type User struct {
	ID        string    `gorm:"primary_key;column:id"` //auto_increment
	Name      string    `gorm:"column:name size:255"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (user *User) Columns() string {
	return "id, name, email, created_at"
}

func (user *User) Fields(u *User) []interface{} {
	return []interface{}{&u.ID, &u.Name, &u.Email, &u.CreatedAt}
}

func (User) TableName() string {
	return "user"
}
