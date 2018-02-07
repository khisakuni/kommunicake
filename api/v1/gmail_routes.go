package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	helpers "github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/jobs/queue"
	"github.com/khisakuni/kommunicake/models"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

func GmailLoginURL(w http.ResponseWriter, r *http.Request) {
	var p oauthLoginURLParams
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}

	if err := p.validate(); err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	redirectURL, _ := url.Parse(p.RedirectURL)
	user := middleware.GetUserFromContext(r.Context())
	authURL, err := formatAuthURL(redirectURL, user, gmailOauthConfig(), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, oauthLoginURLParams{RedirectURL: authURL.String()}, http.StatusOK)
}

func GmailWebhook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	redirectURL, err := url.Parse(query.Get("state"))
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	userID, err := getUserIDFromQuery(query)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	token, err := exchangeCodeForToken(query, func(code string) (*oauth2.Token, error) {
		return gmailOauthConfig().Exchange(oauth2.NoContext, code)
	})
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	_, err = models.FirstOrCreateUserMessageProvider(db, userID, token, models.GMAIL)

	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL.String(), 302)
}

func GmailSubscribeToNewMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := middleware.GetUserFromContext(ctx)
	db := middleware.GetDBFromContext(ctx)
	var provider models.UserMessageProvider
	db.Conn.
		Where(models.UserMessageProvider{UserID: user.ID, MessageProviderType: models.GMAIL}).
		First(&provider)

	token := new(oauth2.Token)
	token.AccessToken = provider.AccessToken
	token.RefreshToken = provider.RefreshToken
	token.Expiry = provider.Expiry
	token.TokenType = provider.TokenType

	config := gmailOauthConfig()
	client := config.Client(ctx, token)
	srv, err := gmail.New(client)
	if err != nil {
		fmt.Println("here!!")
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	req, err := srv.Users.Labels.List("me").Do()
	if err != nil {
		fmt.Printf("OH NO %v\n", err.Error())
	} else {
		fmt.Printf("labels >> %v\n", req.Labels)
	}

	watchReq := gmail.WatchRequest{
		TopicName: "projects/endless-science-125305/topics/new-mail",
		LabelIds:  []string{"INBOX"},
	}
	res, err := srv.Users.Watch("me", &watchReq).Do()
	if err != nil {
		fmt.Println("here 2!!")
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Printf("THINK IT WORKED??? %v\n", res)

	fmt.Fprintf(w, "cool.")
}

func GmailTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := middleware.GetUserFromContext(ctx)

	args := struct {
		Name      string
		HistoryId uint64
		UserId    int
	}{
		Name:      "ProcessGmailHistoryId",
		HistoryId: 3134935,
		UserId:    user.ID,
	}

	q, err := queue.NewQueue("default")
	defer q.CleanUp()
	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(args)
	err = q.Publish(j)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "cool.")
}

func GmailWebhookNewMessage(w http.ResponseWriter, r *http.Request) {
	type PushMessage struct {
		Message struct {
			Data      string `json:"data"`
			MessageID string `json:"message_id"`
		}
		Subscription string `json:"subscription"`
	}
	var message PushMessage
	decodeJSON(w, r.Body, &message)
	fmt.Printf(">>> message >>> data: %s, id: %s, sub: %s\n", message.Message.Data, message.Message.MessageID, message.Subscription)

	fmt.Fprintf(w, "cool")
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
