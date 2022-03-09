package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"time"
)

var (
	hmacSecret []byte
	lgr        logger.Logger
)

func init() {
	cfg, _ := config.NewDefaultConfig()
	hmacSecret = []byte(cfg.AppSecret)
	lgr = logger.NewDefaultLogger()
}

// Claims by token.
type Claims struct {
	UserID    int
	ExpiresAt int64
}

// GetClaimsByToken returns userID by token.
func GetClaimsByToken(tokenString string) (Claims, error) {
	tokenClaims := Claims{}
	defer func() {
		if r := recover(); r != nil {
			lgr.Errorf("%v", r)
		}
	}()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return tokenClaims, fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return tokenClaims, fmt.Errorf("invalid token: %s", tokenString)
	}

	localClaims, ok := claims["token"].(map[string]interface{})
	if !ok {
		return tokenClaims, fmt.Errorf("invalid token claims: %s", tokenString)
	}

	userID, ok := localClaims["UserID"].(float64)
	if !ok {
		return tokenClaims, fmt.Errorf("invalid user id: %v", localClaims["UserID"])
	}

	expiresAt, ok := localClaims["ExpiresAt"].(float64)
	if !ok {
		return tokenClaims, fmt.Errorf("invalid expires at: %v", expiresAt)
	}

	tokenClaims.ExpiresAt = int64(expiresAt)
	tokenClaims.UserID = int(userID)

	return tokenClaims, nil
}

// CreateTokenByUserID returns token string by user.
func CreateTokenByUserID(userID int) (string, error) {
	claims := Claims{
		UserID:    userID,
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": claims,
	})

	return token.SignedString(hmacSecret)
}
