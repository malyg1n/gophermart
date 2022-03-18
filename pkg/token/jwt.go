package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gophermart/pkg/config"
	"time"
)

var (
	hmacSecret []byte
)

func init() {
	cfg, _ := config.NewDefaultConfig()
	hmacSecret = []byte(cfg.AppSecret)
}

// Claims by token.
type Claims struct {
	UserID    uint64
	ExpiresAt int64
}

// GetClaimsByToken returns userID by token.
func GetClaimsByToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenString)
	}

	localClaims, ok := claims["token"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid token claims: %s", tokenString)
	}

	userID, ok := localClaims["UserID"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user id: %v", localClaims["UserID"])
	}

	expiresAt, ok := localClaims["ExpiresAt"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expires at: %v", expiresAt)
	}
	tokenClaims := &Claims{
		ExpiresAt: int64(expiresAt),
		UserID:    uint64(userID),
	}

	return tokenClaims, nil
}

// CreateTokenByUserID returns token string by user.
func CreateTokenByUserID(userID uint64) (string, error) {
	claims := Claims{
		UserID:    userID,
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": claims,
	})

	return token.SignedString(hmacSecret)
}
