package handler

import (
	"encoding/json"
	"errors"
	"gophermart/api/rest/request"
	"gophermart/api/rest/response"
	"gophermart/pkg/contexts"
	"gophermart/pkg/errs"
	"net/http"
)

func (h Handler) ShowBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(contexts.ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "failing with user context", http.StatusInternalServerError)
		return
	}

	balance, err := h.userService.ShowBalance(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorf("%v", err)
		return
	}

	result, err := json.Marshal(balance)
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

// Withdraw handler.
func (h Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(contexts.ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "failing with user context", http.StatusInternalServerError)
		return
	}

	dec := json.NewDecoder(r.Body)

	var req request.Withdraw

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorf("%v", err)
		return
	}

	err := h.userService.Withdraw(ctx, userID, req.Order, req.Sum)
	if err != nil {
		if errors.Is(errs.ErrBalanceTooSmall, err) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		if errors.Is(errs.ErrOrderNumber, err) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorf("%v", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Withdrawals handler.
func (h Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(contexts.ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "failing with user context", http.StatusInternalServerError)
		return
	}

	trans, err := h.userService.GetTransactions(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorf("%v", err)
		return
	}

	if len(trans) == 0 {
		http.Error(w, "no content", http.StatusNoContent)
		return
	}

	withdraws := make([]response.Transaction, 0, len(trans))
	for _, t := range trans {
		withdraws = append(withdraws, response.TransactionFromCanonical(t))
	}

	result, err := json.Marshal(withdraws)
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
