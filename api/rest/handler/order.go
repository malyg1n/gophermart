package handler

import (
	"errors"
	"gophermart/pkg/contexts"
	"gophermart/pkg/errs"
	"io"
	"net/http"
)

// CreateOrder handler.
func (h Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var number string

	userID, ok := ctx.Value(contexts.ContextUserKey).(int)
	if !ok {
		http.Error(w, "failing with user context", http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	number = string(b)

	err = h.orderService.CreateOrder(ctx, number, userID)
	if err != nil {
		if errors.Is(errs.ErrOrderNumber, err) {
			http.Error(w, "incorrect order number", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(errs.ErrOrderCreatedByMyself, err) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(errs.ErrOrderExists, err) {
			http.Error(w, "order already uploaded", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
