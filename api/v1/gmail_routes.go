package v1

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/gmail/v1"
)

func GmailLoginURL(w http.ResponseWriter, r *http.Request) {
  // ctx := context.Background()
	
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

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	res := struct{RedirectURL string}{RedirectURL: authURL }

	jsonResponse(w, res, http.StatusOK)

	// NOTE: This is to exchange token from Gmail API
  // var code string
  // if _, err := fmt.Scan(&code); err != nil {
  //   log.Fatalf("Unable to read authorization code %v", err)
  // }

  // tok, err := config.Exchange(oauth2.NoContext, code)
  // if err != nil {
  //   log.Fatalf("Unable to retrieve token from web %v", err)
  // }
  // return tok
}

func GmailWebhook(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.String())
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "cool: %v\n", query)
}

func GmailExchangeCode(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Code string
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