package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"time"
)

var hmacSecret []byte

func init() {
	cfg, _ := config.GetConfig()
	hmacSecret = []byte(cfg.AppSecret)
}

// Claims by token.
type Claims struct {
	UserID    int
	ExpiresAT int64
}

// GetClaimsByToken returns userID by token.
func GetClaimsByToken(tokenString string) (Claims, error) {
	tokenClaims := Claims{}
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Errorf("%v", r)
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

	expiresAt, ok := localClaims["ExpiresAT"].(float64)
	if !ok {
		return tokenClaims, fmt.Errorf("invalid expires at: %v", expiresAt)
	}

	tokenClaims.ExpiresAT = int64(expiresAt)
	tokenClaims.UserID = int(userID)

	return tokenClaims, nil
}

// CreateTokenByUserID returns token string by user.
func CreateTokenByUserID(userID int) (string, error) {
	claims := Claims{
		UserID:    userID,
		ExpiresAT: time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": claims,
	})

	return token.SignedString(hmacSecret)
}
