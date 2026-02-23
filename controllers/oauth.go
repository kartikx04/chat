package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/goombaio/namegenerator"
	"gorm.io/gorm"

	"github.com/kartikx04/chat/database"
	"github.com/kartikx04/chat/models"
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

	// check if email already in database
	_, err2 := userRepo.GetUserByEmail(authStruct.Email)

	if err2 == gorm.ErrRecordNotFound {
		//generate random username
		seed := time.Now().UTC().UnixNano()
		nameGenerator := namegenerator.NewNameGenerator(seed)

		name := nameGenerator.Generate()

		_, err1 := userRepo.CreateUser(authStruct.Id, authStruct.Email, name, authStruct.Picture)
		if err1 != nil {
			return
		}
	} else if err2 != nil {
		log.Fatal(err2)
	}

	session, _ = utils.Store.Get(req, "userSession")

	session.Values = map[any]any{
		"email":   authStruct.Email,
		"picture": authStruct.Picture,
	}

	session.Save(req, res)

	http.Redirect(res, req, "/home", http.StatusSeeOther)
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

// Home renders a HTML page for logged in users
func Home(res http.ResponseWriter, req *http.Request) {
	session, _ := utils.Store.Get(req, "userSession")

	email, ok := session.Values["email"].(string)
	if !ok {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	picture, _ := session.Values["picture"].(string)

	tmpl, err := template.ParseFiles("assets/home.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(res, map[string]string{
		"Email":   email,
		"Picture": picture,
	})
}
