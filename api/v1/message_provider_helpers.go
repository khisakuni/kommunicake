package v1

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/khisakuni/kommunicake/models"
	"golang.org/x/oauth2"
)

type oauthLoginURLParams struct {
	RedirectURL string `json:"redirectURL"`
}

func (p *oauthLoginURLParams) validate() error {
	if len(p.RedirectURL) == 0 {
		return fmt.Errorf("Must provide redirectURL")
	}

	if _, err := url.Parse(p.RedirectURL); err != nil {
		return err
	}
	return nil
}

func formatAuthURL(redirectURL *url.URL, user *models.User, oauthConfig *oauth2.Config, opts ...oauth2.AuthCodeOption) (*url.URL, error) {
	query := redirectURL.Query()
	query.Add("user_id", strconv.Itoa(user.ID))
	redirectURL.RawQuery = query.Encode()

	stateOption := oauth2.SetAuthURLParam("state", redirectURL.String())
	return url.Parse(oauthConfig.AuthCodeURL("state-token", append(opts, stateOption)...))
}

func exchangeCodeForToken(query url.Values, exchangeFn func(string) (*oauth2.Token, error)) (*oauth2.Token, error) {
	return exchangeFn(query.Get("code"))
}

func getUserIDFromQuery(query url.Values) (int, error) {
	redirectURL, err := url.Parse(query.Get("state"))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(redirectURL.Query().Get("user_id"))
}
