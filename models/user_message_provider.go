package models

import (
	"time"

	"github.com/khisakuni/kommunicake/database"
	"golang.org/x/oauth2"
)

// UserMessageProvider model
type UserMessageProvider struct {
	ID                  int                 `json:"id"`
	UserID              int                 `json:"userId"`
	MessageProviderType MessageProviderType `json:"messageProviderType"`
	RefreshToken        string              `json:"refreshToken"`
	AccessToken         string              `json:"accessToken"`
	Expiry              time.Time
	TokenType           string
}

type MessageProviderType string

const (
	// GMAIL message provider type
	GMAIL MessageProviderType = "GMAIL"

	//k SLACK message provider type
	SLACK MessageProviderType = "SLACK"
)

func FirstOrCreateUserMessageProvider(db *database.DB, userID int, token *oauth2.Token, providerType MessageProviderType) (*UserMessageProvider, error) {
	var provider UserMessageProvider
	db.Conn.
		Where(UserMessageProvider{UserID: userID, MessageProviderType: providerType}).
		Assign(UserMessageProvider{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, Expiry: token.Expiry, TokenType: token.TokenType}).
		FirstOrCreate(&provider)

	return &provider, db.Conn.Error
}
