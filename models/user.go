package models

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

// GenerateToken generates a Token 
func (u *User) GenerateToken(db *database.DB) (*Token, error) {
	type tokenClaims struct {
		Email string
		ID    int
		jwt.StandardClaims
	}
	claims := tokenClaims{
		Email: u.Email,
		ID:    u.ID,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "kommunicake",
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	var t Token
	db.Conn.
		Where(Token{UserID: u.ID}).
		Attrs(Token{UserID: u.ID, CreatedAt: time.Now()}).
		Assign(Token{UpdatedAt: time.Now(), Value: tokenString}).
		FirstOrCreate(&t)

	if db.Conn.Error != nil {
		return nil, db.Conn.Error
	}

	return &t, nil
}
