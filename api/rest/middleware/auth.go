package middleware

import (
	"context"
	"gophermart/pkg/contexts"
	"gophermart/pkg/token"
	"net/http"
	"strings"
	"time"
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

		userClaims, err := token.GetClaimsByToken(tokenString)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if userClaims.ExpiresAt < time.Now().Local().Unix() {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		c := context.WithValue(r.Context(), contexts.ContextUserKey, userClaims.UserID)
		next.ServeHTTP(w, r.WithContext(c))
	})
}
