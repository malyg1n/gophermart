package handler

import "github.com/go-chi/chi/v5"

func (h Handler) GetRouter() chi.Router {
	router := chi.NewRouter()

	router.Post("/api/user/register", h.Register)
	router.Post("/api/user/login", h.Login)

	return router
}
