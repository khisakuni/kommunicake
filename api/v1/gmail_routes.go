package v1

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"

	helpers "github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/models"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

func GmailLoginURL(w http.ResponseWriter, r *http.Request) {
	var p struct {
		RedirectURL string `json:"redirectURL"`
	}
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}

	if len(p.RedirectURL) == 0 {
		helpers.ErrorResponse(w, errors.New("Must provide redirectURL"), http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r.Context())

	redirectURL, err := url.Parse(p.RedirectURL)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	query := redirectURL.Query()
	query.Add("user_id", strconv.Itoa(user.ID))
	redirectURL.RawQuery = query.Encode()

	stateOption := oauth2.SetAuthURLParam("state", redirectURL.String())
	authURL := gmailOauthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline, stateOption, oauth2.ApprovalForce)

	jsonResponse(w, struct{ RedirectURL string }{RedirectURL: authURL}, http.StatusOK)
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

	code := query.Get("code")
	token, err := gmailOauthConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	var provider models.UserMessageProvider
	db.Conn.
		Where(models.UserMessageProvider{UserID: userID, MessageProviderType: "GMAIL"}).
		Assign(models.UserMessageProvider{RefreshToken: token.RefreshToken}).
		FirstOrCreate(&provider)

	if db.Conn.Error != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectUrl.String(), 302)
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
