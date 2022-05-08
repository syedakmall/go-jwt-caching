package route

import (
	"github.com/go-chi/chi/v5"
	h "github.com/syedakmall/malhttp/handler"
	m "github.com/syedakmall/malhttp/middlewares"
)

func UserRoute(r chi.Router) {
	r.Use(m.JWT)
	r.Get("/all", h.UserHandler.GetAllUsers)
	r.Delete("/delete/{id}", h.UserHandler.DeleteUserById)
	r.Get("/home", h.UserHandler.GetUser)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(m.VerifyRedisCache)
		r.Get("/", h.UserHandler.GetUserById)
	})
}
