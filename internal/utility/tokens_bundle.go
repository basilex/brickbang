package utility

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"brickbang/internal/config"
)

type AccessTokenClaims struct {
	UserID    string `json:"uid"`
	JTI       string `json:"jti"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func init() {
	jwtSecret = []byte(config.Get().JWTSecret)
}

func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateAccessToken generates an access token and returns (tokenString, jti, error)
func GenerateAccessToken(userID, sessionID string, ttl time.Duration) (string, string, error) {
	jti := uuid.NewString()
	now := time.Now().UTC()
	claims := AccessTokenClaims{
		UserID:    userID,
		JTI:       jti,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return signed, jti, nil
}

// ParseAccessToken parses token string and returns claims
func ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	tkn, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Ensure HS256
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := tkn.Claims.(*AccessTokenClaims); ok && tkn.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// GenerateRefreshToken creates a refresh token plaintext and returns (plaintext, jti, err)
// format: "<jti>|<random>"
func GenerateRefreshToken(randStr string) (string, string, error) {
	// randStr — может быть utility.GenerateRandomString(32) или т.п.
	if randStr == "" {
		return "", "", errors.New("randStr required")
	}

	jti := uuid.NewString()
	plain := fmt.Sprintf("%s|%s", jti, randStr)
	return plain, jti, nil
}

// ParseRefreshTokenPlain returns jti from plain refresh token "<jti>|<rand>"
func ParseRefreshTokenPlain(plain string) (string, error) {
	parts := strings.SplitN(plain, "|", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", errors.New("invalid refresh token format")
	}
	return parts[0], nil
}
