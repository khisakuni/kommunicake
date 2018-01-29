package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/khisakuni/kommunicake/api/helpers"
	middleware "github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/models"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
)

func SlackLoginURL(w http.ResponseWriter, r *http.Request) {
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
	authURL, err := formatAuthURL(redirectURL, user, slackOAuthConfig())
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, oauthLoginURLParams{RedirectURL: authURL.String()}, http.StatusOK)
}

func SlackWebhook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	redirectURL, err := url.Parse(query.Get("state"))
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Redirect URL > %s\n", redirectURL)

	userID, err := getUserIDFromQuery(query)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	token, err := exchangeCodeForToken(query, func(token string) (*oauth2.Token, error) {
		accessToken, _, err := slack.GetOAuthToken(
			os.Getenv("SLACK_CLIENT_ID"),
			os.Getenv("SLACK_CLIENT_SECRET"),
			query.Get("code"),
			os.Getenv("SLACK_REDIRECT_URL"),
			false,
		)
		if err != nil {
			return nil, err
		}
		return &oauth2.Token{AccessToken: accessToken}, nil
	})
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	_, err = models.FirstOrCreateUserMessageProvider(db, userID, token, models.SLACK)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL.String(), 302)
}

func slackOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("SLACK_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SLACK_REDIRECT_URL"),
		Scopes:       []string{"incoming-webhook"},
		Endpoint: oauth2.Endpoint{
			AuthURL: os.Getenv("SLACK_AUTH_URL"),
		},
	}
}
