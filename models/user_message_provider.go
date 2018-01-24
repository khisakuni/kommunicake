package models

// UserMessageProvider model
type UserMessageProvider struct {
	ID                  int    `json:"id"`
	UserID              int    `json:"userId"`
	MessageProviderType string `json:"messageProviderType"`
	RefreshToken        string `json:"refreshToken"`
}
