package utility

import (
	"brickbang/internal/config"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = config.Get().JWTSecret

// GenerateRandomString — безопасный рандом для refresh токенов
func GenerateRandomString(length int) (string, error) {
	// For hex encoding, we need length/2 bytes (since hex doubles the size)
	byteLen := (length + 1) / 2
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(b)
	return encoded[:length], nil
}

// GenerateAccessToken — создаёт JWT для access токена
func GenerateAccessToken(userID, sessionID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"session": sessionID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseAccessToken — проверяет JWT и возвращает claims
func ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method to prevent algorithm confusion attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
