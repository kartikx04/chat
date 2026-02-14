package routes

import (
	"net/http"

	"github.com/kartikx04/chat/controllers"
)

func UIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", controllers.RenderPage)
	mux.HandleFunc("/home", controllers.Home)
}
