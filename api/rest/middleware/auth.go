package middleware

import (
	"context"
	"gophermart/pkg/contexts"
	"gophermart/pkg/logger"
	"gophermart/pkg/token"
	"net/http"
	"strings"
)

// Auth middleware.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(header, "Bearer ", "", 1)
		if tokenString == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := token.GetUserIDByToken(tokenString)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			logger.GetLogger().Error(err)
			return
		}

		c := context.WithValue(r.Context(), contexts.ContextUserKey, userID)
		next.ServeHTTP(w, r.WithContext(c))
	})
}
