package models

import (
	"time"
	"github.com/khisakuni/kommunicake/database"
)

type Token struct {
	Value string
	UserID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Token) CreateOrUpdate(db *database.DB) error {
	return db.Conn.FirstOrCreate(t).Update(t).Error
}

func (t *Token) Create(db *database.DB) error {
	return db.Conn.Create(t).Error
}

func (t *Token) Update(db *database.DB) error {
	return db.Conn.Update(t).Error
}