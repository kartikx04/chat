package controller

import (
	"github.com/kartikx04/chat/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	utils.OAuthgolang = &oauth2.Config{
		RedirectURL:  "",
		ClientID:     "",
		ClientSecret: "",
		// scopes limits the access given to a token. this scope returns just the user info of the
		// signed in email address
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint, //Endpoint is Google's OAuth 2.0 default endpoint
	}
}
