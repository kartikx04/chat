package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/models"
	redisrepo "github.com/kartikx04/chat/internal/redis-repo"
	"github.com/kartikx04/chat/pkg"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	pkg.OAuthgolang = &oauth2.Config{
		RedirectURL:  pkg.LoadFile("REDIRECT_URL"),
		ClientID:     pkg.LoadFile("CLIENT_ID"),
		ClientSecret: pkg.LoadFile("CLIENT_SECRET"),
		// scopes limits the access given to a token. this scope returns just the user info of the
		// signed in email address
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint, //Endpoint is Google's OAuth 2.0 default endpoint
	}

}

// Through GoogleSignOn, a URL is returned to the consent page created on Google Console. A security token string is provided which will be parsed and verified during redirect callback.
func GoogleSignOn(res http.ResponseWriter, req *http.Request) {
	tokenString, err := pkg.TokenString()
	if err != nil {
		fmt.Fprintf(res, "error: could not generate random token string: %v", err)
	}

	// creates a new session
	session, err := pkg.Store.Get(req, "tokenSession")
	if err != nil {
		fmt.Fprintf(res, "error: %v", err)
	}

	// saves the generated token string into the created session; uses tokenStringKey as the key
	session.Values["tokenStringKey"] = tokenString
	session.Save(req, res)

	// returns a URL with attached tokenString
	url := pkg.OAuthgolang.AuthCodeURL(tokenString)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// Callback is triggered when Google Cloud Console redirects to /callback. Coupled with GetUserData
// it verifies authorization request and unmarshals returned user data into our created OAuthData data structure
func Callback(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	code := req.FormValue("code")

	// returns the created session
	session, err := pkg.Store.Get(req, "tokenSession")
	if err != nil {
		fmt.Fprintf(res, "error: %v", err)
	}

	// returns the value of tokenStringKey
	dataToken, ok := session.Values["tokenStringKey"].(string)
	if !ok {
		http.Error(res, "session expired or invalid", 400)
		return
	}

	data, err := pkg.GetUserData(state, code, dataToken)
	if err != nil {
		log.Println("GetUserData error:", err)
		http.Error(res, err.Error(), 500)
		return
	}

	// the session cookie is deleted immediately
	session.Options.MaxAge = -1
	session.Save(req, res)

	var authStruct models.OAuthData

	// Google Cloud Console returns a JSON structure containing "id",,"email", "verified_email" and "picture"
	// this converts the JSON structure into our created OAuthData structure
	err = json.Unmarshal([]byte(data), &authStruct)
	if err != nil {
		fmt.Fprintf(res, "error: %v", err)
	}

	// if the email is valid then add the user information to cookie and save it.
	if !authStruct.VerifiedEmail {
		return
	}

	userRepo := database.NewUserRepository(database.DB)

	seed := time.Now().UTC().UnixNano()
	name := namegenerator.NewNameGenerator(seed).Generate()

	user, err := userRepo.GetOrCreateUser(authStruct.Id, authStruct.Email, name, authStruct.Picture)
	if err != nil {
		log.Println("GetorCreateUser error:", err)
		http.Error(res, err.Error(), 500)
		return
	}

	// Set Redis lookup keys so frontend can resolve username <-> id
	redisrepo.SetUsernameLookup(user.Id, user.Username)
	redisrepo.SetIdLookup(user.Username, user.Id)

	session, _ = pkg.Store.Get(req, "userSession")

	session.Values = map[any]any{
		"email":   authStruct.Email,
		"picture": authStruct.Picture,
	}

	session.Save(req, res)

	log.Printf("OAuth Callback: user.Id=%s, user.Username=%s", user.Id.String(), user.Username)

	frontendURL := pkg.LoadFile("FRONTEND_URL")

	redirectURL := fmt.Sprintf(
		"%s/auth/callback?id=%s&username=%s&email=%s",
		frontendURL,
		user.Id.String(),
		url.QueryEscape(user.Username),
		url.QueryEscape(authStruct.Email),
	)
	http.Redirect(res, req, redirectURL, http.StatusFound)

}

func Logout(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:     "userSession",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Set MaxAge to -1 to delete the cookie
		HttpOnly: true,
	})
	http.Redirect(res, req, "/", http.StatusSeeOther)
	fmt.Printf("user logged out successfully")
}
