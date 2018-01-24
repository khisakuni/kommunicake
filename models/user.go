package models

import (
	"github.com/khisakuni/kommunicake/database"
)

// User model
type User struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
}

// Create writes user data to db
func (u *User) Create(db *database.DB) error {
	return db.Conn.Create(u).Error
}
