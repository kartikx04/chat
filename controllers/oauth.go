package controller

import (
	"fmt"
	"net/http"

	"github.com/kartikx04/chat/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	utils.OAuthgolang = &oauth2.Config{
		RedirectURL:  utils.LoadFile("REDIRECT_URL"),
		ClientID:     utils.LoadFile("CLIENT_ID"),
		ClientSecret: utils.LoadFile("CLIENT_SECRET"),
		// scopes limits the access given to a token. this scope returns just the user info of the
		// signed in email address
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint, //Endpoint is Google's OAuth 2.0 default endpoint
	}
}

// Through GoogleSignOn, a URL is returned to the consent page created on Google Console. A security token string is provided which will be parsed and verified during redirect callback.
func GoogleSignOn(res http.ResponseWriter, req *http.Request) {
	tokenString, err := utils.TokenString()
	if err != nil {
		fmt.Fprintf(res, "error: could not generate random token string: %v", err)
	}

	// creates a new session
	session, err := utils.Store.Get(req, "tokenSession")
	if err != nil {
		fmt.Fprintf(res, "error: %v", err)
	}

	// saves the generated token string into the created session; uses tokenStringKey as the key
	session.Values["tokenStringKey"] = tokenString
	session.Save(req, res)

	// returns a URL with attached tokenString
	url := utils.OAuthgolang.AuthCodeURL(tokenString)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
