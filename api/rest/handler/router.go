package handler

import (
	"github.com/go-chi/chi/v5"
	"gophermart/api/rest/middleware"
)

// GetRouter returns router.
func (h Handler) GetRouter() chi.Router {
	router := chi.NewRouter().With(
		middleware.Compress,
		middleware.Decompress,
	)

	router.Post("/api/user/register", h.Register)
	router.Post("/api/user/login", h.Login)

	router.Route("/api/user/orders", func(r chi.Router) {
		r = r.With(middleware.Auth)
		r.Post("/", h.CreateOrder)
		r.Get("/", h.GetOrdersByUser)
	})

	router.Route("/api/user/balance", func(r chi.Router) {
		r = r.With(middleware.Auth)
	})

	return router
}
