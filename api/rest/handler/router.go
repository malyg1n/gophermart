package handler

import (
	"github.com/go-chi/chi/v5"
	"gophermart/api/rest/middleware"
	"net/http"
)

// GetRouter returns router.
func (h Handler) GetRouter() chi.Router {
	router := chi.NewRouter().With(
		middleware.Compress,
		middleware.Decompress,
	)

	router.Post("/api/user/register", h.Register)
	router.Post("/api/user/login", h.Login)

	router.Route("/api/order", func(r chi.Router) {
		r = r.With(middleware.Auth)
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(202)
		})
	})

	return router
}
