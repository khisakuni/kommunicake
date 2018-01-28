package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	helpers "github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/models"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type gmailLoginURLParams struct {
	RedirectURL string `json:"redirectURL"`
}

func (p *gmailLoginURLParams) validate() error {
	if len(p.RedirectURL) == 0 {
		return fmt.Errorf("Must provide redirectURL")
	}

	if _, err := url.Parse(p.RedirectURL); err != nil {
		return err
	}
	return nil
}

func GmailLoginURL(w http.ResponseWriter, r *http.Request) {
	var p gmailLoginURLParams
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}

	if err := p.validate(); err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	redirectURL, _ := url.Parse(p.RedirectURL)
	user := middleware.GetUserFromContext(r.Context())
	authURL := formatGmailAuthURL(redirectURL, user)

	jsonResponse(w, gmailLoginURLParams{RedirectURL: authURL}, http.StatusOK)
}

func GmailWebhook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	redirectUrl, err := url.Parse(query.Get("state"))
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	redirectQuery := redirectUrl.Query()
	userID, err := strconv.Atoi(redirectQuery.Get("user_id"))
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	token, err := exchangeCodeForToken(query.Get("code"))
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	err = firstOrCreateUserMessageProvider(db, userID, token)

	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectUrl.String(), 302)
}

func firstOrCreateUserMessageProvider(db *database.DB, userID int, token *oauth2.Token) error {
	var provider models.UserMessageProvider
	db.Conn.
		Where(models.UserMessageProvider{UserID: userID, MessageProviderType: "GMAIL"}).
		Assign(models.UserMessageProvider{RefreshToken: token.RefreshToken}).
		FirstOrCreate(&provider)

	return db.Conn.Error
}

func exchangeCodeForToken(code string) (*oauth2.Token, error) {
	return gmailOauthConfig().Exchange(oauth2.NoContext, code)
}

func formatGmailAuthURL(redirectURL *url.URL, user *models.User) string {
	query := redirectURL.Query()
	query.Add("user_id", strconv.Itoa(user.ID))
	redirectURL.RawQuery = query.Encode()

	stateOption := oauth2.SetAuthURLParam("state", redirectURL.String())
	return gmailOauthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline, stateOption, oauth2.ApprovalForce)
}

func gmailOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GMAIL_REDIRECT_URL"),
		Scopes:       []string{gmail.GmailComposeScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("GMAIL_AUTH_URL"),
			TokenURL: os.Getenv("GMAIL_TOKEN_URL"),
		},
	}
}
