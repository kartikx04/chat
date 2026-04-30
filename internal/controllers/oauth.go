package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/kartikx04/chat/internal/auth"
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
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleSignOn(res http.ResponseWriter, req *http.Request) {
	tokenString, err := pkg.TokenString()
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to generate oauth state token", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	session, err := pkg.Store.Get(req, "tokenSession")
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to get token session", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["tokenStringKey"] = tokenString
	session.Save(req, res)

	authURL := pkg.OAuthgolang.AuthCodeURL(tokenString)
	slog.InfoContext(req.Context(), "oauth redirect initiated")
	http.Redirect(res, req, authURL, http.StatusTemporaryRedirect)
}

func Callback(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	code := req.FormValue("code")

	session, err := pkg.Store.Get(req, "tokenSession")
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to get token session", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	dataToken, ok := session.Values["tokenStringKey"].(string)
	if !ok {
		slog.WarnContext(req.Context(), "oauth callback: session expired or invalid")
		http.Error(res, "session expired or invalid", http.StatusBadRequest)
		return
	}

	data, err := pkg.GetUserData(state, code, dataToken)
	if err != nil {
		slog.ErrorContext(req.Context(), "oauth get user data failed", "error", err)
		http.Error(res, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Delete token session immediately after use
	session.Options.MaxAge = -1
	session.Save(req, res)

	var authStruct models.OAuthData
	if err := json.Unmarshal([]byte(data), &authStruct); err != nil {
		slog.ErrorContext(req.Context(), "failed to unmarshal oauth data", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	if !authStruct.VerifiedEmail {
		slog.WarnContext(req.Context(), "oauth callback: unverified email rejected", "email", authStruct.Email)
		http.Error(res, "email not verified", http.StatusForbidden)
		return
	}

	userRepo := database.NewUserRepository(database.DB)
	seed := time.Now().UTC().UnixNano()
	name := namegenerator.NewNameGenerator(seed).Generate()

	user, err := userRepo.GetOrCreateUser(authStruct.Id, authStruct.Email, name, authStruct.Picture)
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to get or create user", "error", err, "email", authStruct.Email)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	redisrepo.SetUsernameLookup(user.Id, user.Username)
	redisrepo.SetIdLookup(user.Username, user.Id)

	userSession, _ := pkg.Store.Get(req, "userSession")
	userSession.Values = map[any]any{
		"email":   authStruct.Email,
		"picture": authStruct.Picture,
	}
	userSession.Save(req, res)

	slog.InfoContext(req.Context(), "oauth login success",
		"user_id", user.Id.String(),
		"username", user.Username,
	)

	token, err := auth.GenerateToken(user.Id.String(), user.Username, authStruct.Email)
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to generate jwt", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	frontendURL := pkg.LoadFile("FRONTEND_URL")
	http.Redirect(res, req, frontendURL+"/auth/callback?token="+token, http.StatusFound)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:     "userSession",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	slog.InfoContext(req.Context(), "user logged out")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// internal/controllers/auth.go
func Me(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		slog.WarnContext(req.Context(), "me: no auth header")
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := auth.ValidateToken(tokenStr)
	if err != nil {
		slog.WarnContext(req.Context(), "me: invalid token", "error", err)
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	slog.DebugContext(req.Context(), "me: identity resolved",
		"user_id", claims.UserID,
		"username", claims.Username,
	)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{
		"id":       claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
	})
}
