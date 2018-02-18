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

type params struct {
	db        *database.DB
	userID    int
	historyID uint64
	service   *gmail.Service
	provider  *models.UserMessageProvider
}

// ProcessGmailHistoryID does something idk lol
func ProcessGmailHistoryID(historyID uint64, userID int) {
	fmt.Printf("WOO HOO %d, %d\n", historyID, userID)
	//	fmt.Println("Starting...")
	//	db, err := database.NewDB()
	//	handleError(err)
	//	defer db.Conn.Close()
	//
	//	p := &params{db: db, userID: userID, historyID: historyID}
	//
	//	fmt.Println("Gmail service...")
	//	err = gmailService(p)
	//	handleError(err)
	//	fmt.Println("Done.")
	//
	//	fmt.Println("Processing messages...")
	//	err = processMessages(p)
	//	handleError(err)
	//	fmt.Println("Done.")
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Uh oh >> %s\n", err.Error())
		panic(err.Error())
	}
}

func processMessages(p *params) error {
	if p.provider.HistoryID > 0 {
		// partial sync
		res, err := p.service.Users.History.List("me").StartHistoryId(p.provider.HistoryID).Do()
		if err != nil {
			return err
		}

		p.provider.HistoryID = p.historyID
		p.db.Conn.Model(p.provider).Update("history_id")
		if p.db.Conn.Error != nil {
			return p.db.Conn.Error
		}

		fmt.Printf("HISTORY COUNT %d\n", len(res.History))
		for _, history := range res.History {
			fmt.Printf("MESSSAGES ADDED COUNT %d\n", len(history.MessagesAdded))
			for _, m := range history.MessagesAdded {
				err = processMessage(p.service, m.Message.Id)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	// "full" sync
	fmt.Println("FULL SYNC")

	p.provider.HistoryID = p.historyID
	p.db.Conn.Model(p.provider).Update("history_id")
	if p.db.Conn.Error != nil {
		return p.db.Conn.Error
	}

	return nil
}

func processMessage(srv *gmail.Service, messageID string) error {
	messageRes, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return err
	}

	fmt.Printf(">>> DATA >>> %s\n", messageRes.Payload.Body.Data)
	fmt.Printf("PARTS COUNT: %d\n", len(messageRes.Payload.Parts))
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

func gmailService(p *params) error {
	provider, err := provider(p.db, p.userID)
	if err != nil {
		return err
	}
	p.provider = provider

	token, err := tokenFromUser(provider)
	if err != nil {
		return err
	}

	config := gmailOauthConfig()
	ctx := context.Background()
	client := config.Client(ctx, token)
	srv, err := gmail.New(client)
	if err != nil {
		return err
	}
	p.service = srv
	return nil
}

func tokenFromUser(provider *models.UserMessageProvider) (*oauth2.Token, error) {
	token := new(oauth2.Token)
	token.AccessToken = provider.AccessToken
	token.RefreshToken = provider.RefreshToken
	token.Expiry = provider.Expiry
	return token, nil
}

func provider(db *database.DB, userID int) (*models.UserMessageProvider, error) {
	var provider models.UserMessageProvider
	db.Conn.
		Where(models.UserMessageProvider{UserID: userID, MessageProviderType: models.GMAIL}).
		First(&provider)
	return &provider, db.Conn.Error
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
