package v1

import (
	"errors"
	"os"
	"net/http"
	"net/url"
  "golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	helpers "github.com/khisakuni/kommunicake/api/helpers"
)

func GmailLoginURL(w http.ResponseWriter, r *http.Request) {
	var p struct {RedirectURL string `json:"redirectURL"`}
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}

	if len(p.RedirectURL) == 0 {
		helpers.ErrorResponse(w, errors.New("Must provide redirectURL"), http.StatusBadRequest)
		return
	}
	
	stateOption := oauth2.SetAuthURLParam("state", p.RedirectURL)
	authURL := gmailOauthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline, stateOption)

	jsonResponse(w, struct{RedirectURL string}{RedirectURL: authURL}, http.StatusOK)
}

func GmailWebhook(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.String())
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	url, ok := query["/api/v1/webhooks/gmail?state"]
	if !ok {
		helpers.ErrorResponse(w, errors.New("Need redirect url"), http.StatusBadRequest)
		return
	}

	// TODO: Store refresh token and auth token in DB and use to make requests to GMAIL API
	// code, ok := query["code"]
	// if !ok {
	// 	ErrorResponse(w, errors.New("Missing code"), http.StatusInternalServerError)
	// 	return
	// }

	// tok, err := gmailOauthConfig().Exchange(oauth2.NoContext, code[0])
  // if err != nil {
  //   log.Fatalf("Unable to retrieve token from web %v", err)
	// }

	http.Redirect(w, r, url[0], 302)
}

func gmailOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		RedirectURL: os.Getenv("GMAIL_REDIRECT_URL"),
		Scopes: []string{gmail.GmailComposeScope},
		Endpoint: oauth2.Endpoint{
			AuthURL: os.Getenv("GMAIL_AUTH_URL"),
			TokenURL: os.Getenv("GMAIL_TOKEN_URL"),
		},
	}
}