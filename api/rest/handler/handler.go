package handler

import "gophermart/service"

// Handler base handler struct.
type Handler struct {
	userService service.IUserService
}

type Option func(handler *Handler)

// NewHandler returns Handler instance.
func NewHandler(opts ...Option) *Handler {
	handler := &Handler{}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

// WithUserService option.
func WithUserService(sv service.IUserService) Option {
	return func(handler *Handler) {
		handler.userService = sv
	}
}
