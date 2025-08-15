package jwt

import (
	"github.com/golang-jwt/jwt"
	"serversTest2/internal/config"
	"time"
)

func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(config.JwtKey)
}
