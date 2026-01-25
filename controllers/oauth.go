package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

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

// Callback is triggered when Google Cloud Console redirects to /callback. Coupled with GetUserData
// it verifies authorization request and unmarshals returned user data into our created OAuthData data structure
func Callback(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	code := req.FormValue("code")

	// returns the created session
	session, err := utils.Store.Get(req, "tokenSession")
	if err != nil {
		fmt.Fprintf(res, "error: %v", err)
	}

	// returns the value of tokenStringKey
	dataToken, ok := session.Values["tokenStringKey"].(string)
	if !ok {
		dataToken = "token not found in the session"
	}

	data, err := utils.GetUserData(state, code, dataToken)
	if err != nil {
		log.Fatal(err)
	}

	// the session cookie is deleted immediately
	session.Options.MaxAge = -1
	session.Save(req, res)

	var authStruct utils.OAuthData

	// Google Cloud Console returns a JSON structure containing "id",,"email", "verified_email" and "picture"
	// this converts the JSON structure into our created OAuthData structure
	err = json.Unmarshal([]byte(data), &authStruct)
	if err != nil {
		fmt.Fprintf(res, "error: %v", err)
	}

	// returns a response a response based on verification success or failure
	status := authStruct.Verified_email
	if status {
		fmt.Fprintf(res, "success: %s is a verified user\n", authStruct.Email)
	} else {
		fmt.Fprint(res, "failed verification")
	}
}

// RenderPage renders a simple HTML page to try out Google Sign-On
func RenderPage(res http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("assets/index.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(res, nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
