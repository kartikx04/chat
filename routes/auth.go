package routes

import (
	"net/http"

	"github.com/kartikx04/chat/controllers"
)

func AuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/google-sso", controllers.GoogleSignOn)
	mux.HandleFunc("/callback", controllers.Callback)
}
