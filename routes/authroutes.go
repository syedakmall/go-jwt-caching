package route

import (
	"github.com/go-chi/chi/v5"
	h "github.com/syedakmall/malhttp/handler"
)

func AuthRoute(r chi.Router) {
	r.Post("/signin", h.AuthHandler.SignIn)
	r.Post("/signup", h.AuthHandler.SignUp)
	r.Get("/refresh", h.AuthHandler.RefreshToken)
}
