package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/jobs/queue"
	"github.com/khisakuni/kommunicake/models"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	q, err := queue.NewQueue("default")
	defer q.CleanUp()
	handleError(err)

	q.RegisterWorker("SayHi", func(args interface{}) {
		bodyMap := args.(map[string]interface{})
		message := bodyMap["Message"].(string)
		fmt.Printf("hi!! %s\n", message)
	})

	q.RegisterWorker("ProcessGmailHistoryId", func(args interface{}) {
		fmt.Println("received message!")
		bodyMap := args.(map[string]interface{})
		fmt.Println(bodyMap)
		historyId := uint64(bodyMap["HistoryId"].(float64))
		userId := int(bodyMap["UserId"].(float64))

		db, err := database.NewDB()
		defer db.Conn.Close()
		if err != nil {
			panic(err)
		}

		var provider models.UserMessageProvider
		db.Conn.
			Where(models.UserMessageProvider{UserID: userId, MessageProviderType: models.GMAIL}).
			First(&provider)

		token := new(oauth2.Token)
		token.AccessToken = provider.AccessToken
		token.RefreshToken = provider.RefreshToken
		token.Expiry = provider.Expiry
		token.TokenType = provider.TokenType

		config := gmailOauthConfig()
		ctx := context.Background()
		client := config.Client(ctx, token)
		srv, err := gmail.New(client)

		if err != nil {
			panic(err)
		}

		res, err := srv.Users.History.List("me").StartHistoryId(historyId).Do()
		if err != nil {
			panic(err.Error())
		}

		h := res.History[0]

		for i, m := range h.MessagesAdded {
			fmt.Printf("snippet > %d: %s\n", i, m.Message.Id)
			messageRes, err := srv.Users.Messages.Get("me", m.Message.Id).Do()
			if err != nil {
				panic(err)
			}

			for _, part := range messageRes.Payload.Parts {
				fmt.Printf(">>>> mime: %s\n", part.MimeType)
			}
		}

	})

	q.Consume()

	select {}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
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
