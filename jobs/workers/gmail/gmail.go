package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/models"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
)

func ProcessGmailHistoryId(historyID uint64, userID int) {
	db, err := database.NewDB()
	defer db.Conn.Close()
	handleError(err)

	srv, err := gmailService(db, userID)
	handleError(err)

	err = processMessages(srv, historyID)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Uh oh >> %s\n", err.Error())
		panic(err.Error())
	}
}

func processMessages(srv *gmail.Service, historyID uint64) error {
	res, err := srv.Users.History.List("me").StartHistoryId(historyID).Do()
	if err != nil {
		return err
	}

	for _, history := range res.History {
		for _, m := range history.MessagesAdded {
			err = processMessage(srv, m.Message.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processMessage(srv *gmail.Service, messageID string) error {
	messageRes, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return err
	}

	for _, part := range messageRes.Payload.Parts {
		fmt.Printf(">>>> mime: %s\n", part.MimeType)
		decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
		if err != nil {
			return err
		}
		fmt.Printf("message: %s\n", string(decoded))

	}

	return nil
}

func gmailService(db *database.DB, userID int) (*gmail.Service, error) {
	token, err := tokenFromUser(db, userID)
	if err != nil {
		return nil, err
	}

	config := gmailOauthConfig()
	ctx := context.Background()
	client := config.Client(ctx, token)
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func tokenFromUser(db *database.DB, userID int) (*oauth2.Token, error) {
	var provider models.UserMessageProvider
	db.Conn.
		Where(models.UserMessageProvider{UserID: userID, MessageProviderType: models.GMAIL}).
		First(&provider)
	if db.Conn.Error != nil {
		return nil, db.Conn.Error
	}

	token := new(oauth2.Token)
	token.AccessToken = provider.AccessToken
	token.RefreshToken = provider.RefreshToken
	token.Expiry = provider.Expiry
	return token, nil
}

func gmailOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GMAIL_REDIRECT_URL"),
		Scopes:       []string{gmail.GmailComposeScope, gmail.GmailLabelsScope, gmail.GmailModifyScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("GMAIL_AUTH_URL"),
			TokenURL: os.Getenv("GMAIL_TOKEN_URL"),
		},
	}
}
