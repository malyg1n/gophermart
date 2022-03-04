package handler

import (
	"encoding/json"
	"errors"
	"gophermart/api/rest/response"
	"gophermart/pkg/contexts"
	"gophermart/pkg/errs"
	"io"
	"net/http"
)

// CreateOrder handler.
func (h Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	number := string(b)
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
		h.logger.Errorf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetOrdersByUser handler.
func (h Handler) GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(contexts.ContextUserKey).(int)
	if !ok {
		http.Error(w, "failing with user context", http.StatusInternalServerError)
		return
	}

	orders, err := h.orderService.GetOrdersByUser(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorf("%v", err)
		return
	}

	if len(orders) == 0 {
		http.Error(w, "orders did not upload", http.StatusNoContent)
		return
	}

	responseOrders := make([]response.Order, 0, len(orders))
	for _, o := range orders {
		responseOrders = append(responseOrders, response.OrderFromCanonical(o))
	}

	result, err := json.Marshal(responseOrders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorf("%v", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(result)
	if err != nil {
		h.logger.Errorf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
