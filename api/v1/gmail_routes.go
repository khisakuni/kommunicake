package v1

import (
	"fmt"
	"errors"
  "io/ioutil"
  "log"
	"net/http"
	"net/url"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/gmail/v1"
)

// 1.) Get redirect URL from request
// 2.) 
func GmailLoginURL(w http.ResponseWriter, r *http.Request) {
	// ctx := context.Background()
	type requestParams struct {
		RedirectURL string `json:"redirectURL`
	}
	var p requestParams
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}
	
	b, err := ioutil.ReadFile("/Users/koheihisakuni/gmail_client_secrets.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("state", p.RedirectURL))

  fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	res := struct{RedirectURL string}{RedirectURL: authURL }

	jsonResponse(w, res, http.StatusOK)
}

func GmailWebhook(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.String())
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	url, ok := query["/api/v1/webhooks/gmail?state"]
	if !ok {
		errorResponse(w, errors.New("Need redirect url"), http.StatusBadRequest)
	}

	code, ok := query["code"]
	if !ok {
		errorResponse(w, errors.New("Missing code"), http.StatusInternalServerError)
	}
	// fmt.Fprintf(w, "url: %s, code: %s\n", url, code)
	
	// vvvvvvvvvvvvvvv GET ACCESS TOKEN vvvvvvvvvvvv
	b, err := ioutil.ReadFile("/Users/koheihisakuni/gmail_client_secrets.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope)
	tok, err := config.Exchange(oauth2.NoContext, code[0])
  if err != nil {
    log.Fatalf("Unable to retrieve token from web %v", err)
	}
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// TODO: Store refresh token and auth token in DB and use to make requests to GMAIL API
	fmt.Println(tok)

	http.Redirect(w, r, url[0], 302)
}

func GmailExchangeCode(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Code string `json:"code"`
	}
	var p params
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}

	b, err := ioutil.ReadFile("/Users/koheihisakuni/gmail_client_secrets.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope)
	tok, err := config.Exchange(oauth2.NoContext, p.Code)
  if err != nil {
    log.Fatalf("Unable to retrieve token from web %v", err)
	}

	type response struct {
		Token *oauth2.Token
	}
	
	jsonResponse(w, response{Token: tok}, http.StatusOK)
}