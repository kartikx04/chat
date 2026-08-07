package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/controller"
)

func Register(r *chi.Mux) {
	r.Get("/health", controller.Health)

	r.Get("/auth/google-sso", controller.GoogleSignOn)
	r.Get("/auth/google/callback", controller.Callback)

	r.Get("/home", controller.Home)
}
