package models

import (
	"fmt"
	"time"
	"os"
	"github.com/khisakuni/kommunicake/database"
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	Value string
	UserID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenerateToken(db *database.DB, user *User) (*Token, error) {
	tokenString, err := newTokenString(user)
	if err != nil {
		return nil, err
	}

	var t Token
	db.Conn.
		Where(Token{UserID: user.ID}).
		Attrs(Token{UserID: user.ID, CreatedAt: time.Now()}).
		Assign(Token{UpdatedAt: time.Now(), Value: tokenString}).
		FirstOrCreate(&t)

	if db.Conn.Error != nil {
		return nil, db.Conn.Error
	}

	return &t, nil
}

func FromTokenString(tokenString string, db *database.DB) (*Token, error) {
	userID, err := userIDFromTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	var token Token
	db.Conn.Where(Token{UserID: int(userID)}).First(&token)
	return &token, db.Conn.Error	
}

func (t *Token) GetUser(db *database.DB) (*User, error) {
	userID, err := userIDFromTokenString(t.Value)
	if err != nil {
		return nil, err
	}

	var user User
	db.Conn.First(&user, userID)
	return &user, db.Conn.Error
}

func userIDFromTokenString(tokenString string) (int, error) {
	claims, err := claimsFromTokenString(tokenString)
	if err != nil {
		return 0, err
	}

	if userID, ok := claims["ID"]; ok {
		if id, ok := userID.(float64); ok {
			return int(id), nil	
		}
	}

	return 0, fmt.Errorf("Invalid token")
}

func claimsFromTokenString(tokenString string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		return nil, ve
	} else {
		return nil, fmt.Errorf("Invalid token")
	}
}

func newTokenString(user *User) (string, error) {
	type tokenClaims struct {
		Email string
		ID    int
		jwt.StandardClaims
	}
	claims := tokenClaims{
		Email: user.Email,
		ID:    user.ID,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "kommunicake",
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
