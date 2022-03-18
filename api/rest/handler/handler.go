package handler

import (
	"gophermart/pkg/logger"
	"gophermart/service"
)

// Handler base handler struct.
type Handler struct {
	userService  service.UserProcessor
	orderService service.OrderProcessor
	logger       logger.Logger
}

// Option for Handler.
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
func WithUserService(sv service.UserProcessor) Option {
	return func(handler *Handler) {
		handler.userService = sv
	}
}

// WithOrderService option.
func WithOrderService(sv service.OrderProcessor) Option {
	return func(handler *Handler) {
		handler.orderService = sv
	}
}

// WithLogger option.
func WithLogger(l logger.Logger) Option {
	return func(handler *Handler) {
		handler.logger = l
	}
}
