package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gophermart/pkg/config"
)

var hmacSecret []byte

func init() {
	cfg, _ := config.GetConfig()
	hmacSecret = []byte(cfg.AppSecret)
}

// GetUserIDByToken returns userID by token.
func GetUserIDByToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token: %s", tokenString)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token claims: %s", tokenString)
	}

	return int(userID), nil
}

// CreateTokenByUserID returns token string by user.
func CreateTokenByUserID(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	return token.SignedString(hmacSecret)
}
