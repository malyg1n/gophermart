package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/api/rest/request"
	"gophermart/pkg/errs"
	"gophermart/pkg/logger"
	"net/http"
)

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dec := json.NewDecoder(r.Body)

	var req request.AuthRequest

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.GetLogger().Error(err)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	err := h.userService.Create(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(errs.ErrLoginExists, err) {
			logger.GetLogger().Info(err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.GetLogger().Error(err)
		return
	}

	h.login(w, ctx, req.Login, req.Password)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dec := json.NewDecoder(r.Body)

	var req request.AuthRequest

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.GetLogger().Error(err)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	h.login(w, ctx, req.Login, req.Password)
}

func (h Handler) login(w http.ResponseWriter, ctx context.Context, login, password string) {
	token, err := h.userService.Auth(ctx, login, password)
	if err != nil {
		if errors.Is(errs.ErrAuthFailed, err) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}
