package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"strings"
	"time"

	"brickbang/internal/exception"
	"brickbang/internal/repository"
	"brickbang/storage/dbs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(ctx *fiber.Ctx) (map[string]any, error)
	Login(ctx *fiber.Ctx) (map[string]any, error)
	Refresh(ctx *fiber.Ctx) (map[string]any, error)
	Me(ctx *fiber.Ctx) (map[string]any, error)

	CreateAPIKey(ctx *fiber.Ctx) (map[string]any, error)
	DeleteAPIKey(ctx *fiber.Ctx) (map[string]any, error)
	ListAPIKeys(ctx *fiber.Ctx) ([]*dbs.Apikey, error)
}

type AuthService struct {
	repo      repository.IAuthRepository
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(repo repository.IAuthRepository, jwtSecret string, jwtTTL time.Duration) IAuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

// ============ Registration ============
func (s *AuthService) Register(ctx *fiber.Ctx) (map[string]any, error) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return nil, exception.ErrBadRequest("invalid JSON body")
	}
	if body.Username == "" || body.Password == "" {
		return nil, exception.ErrUnprocessable("username and password are required")
	}

	existing, _ := s.repo.FindByUsername(context.Background(), body.Username)
	if existing != nil {
		return nil, exception.ErrConflict("user already exists")
	}

	hashed := HashPassword(body.Password)

	user, err := s.repo.CreateUser(context.Background(), &dbs.AuthCreateUserParams{
		Username: body.Username,
		Password: hashed,
	})
	if err != nil {
		return nil, exception.ErrInternal("failed to create user")
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, exception.ErrInternal("failed to generate token")
	}

	return map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
		"token": token,
	}, nil
}

// ============ Login ============
func (s *AuthService) Login(ctx *fiber.Ctx) (map[string]any, error) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return nil, exception.ErrBadRequest("invalid JSON body")
	}

	user, err := s.repo.FindByUsername(context.Background(), body.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrUnauthorized("invalid username or password")
		}
		return nil, exception.ErrInternal("failed to fetch user")
	}
	if user == nil || !CheckPasswordHash(body.Password, user.Password) {
		return nil, exception.ErrUnauthorized("invalid username or password")
	}
	if user.IsBlocked {
		return nil, exception.ErrForbidden("user is blocked")
	}

	_, _ = s.repo.UpdateVisitedAt(context.Background(), user.ID)

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, exception.ErrInternal("failed to generate token")
	}

	return map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
		"token": token,
	}, nil
}

// ============ Refresh ============
func (s *AuthService) Refresh(ctx *fiber.Ctx) (map[string]any, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	newToken, err := s.generateJWT(userID)
	if err != nil {
		return nil, exception.ErrInternal("failed to generate new token")
	}
	return map[string]any{"token": newToken}, nil
}

// ============ Me ============
func (s *AuthService) Me(ctx *fiber.Ctx) (map[string]any, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(context.Background(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNotFound("user not found")
		}
		return nil, exception.ErrInternal("failed to fetch user")
	}

	return map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"is_checked": user.IsChecked,
		"is_blocked": user.IsBlocked,
		"visited_at": user.VisitedAt,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}, nil
}

// ============ API Keys ============
func (s *AuthService) CreateAPIKey(ctx *fiber.Ctx) (map[string]any, error) {
	var body struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return nil, exception.ErrBadRequest("invalid JSON body")
	}
	if body.UserID == "" || body.Name == "" {
		return nil, exception.ErrUnprocessable("missing required fields: user_id, name")
	}

	rawKey := GenerateAPIKey(32)
	keyHash := HashAPIKey(rawKey)

	apikey, err := s.repo.CreateAPIKey(context.Background(), &dbs.AuthCreateAPIKeyParams{
		UserID:  body.UserID,
		KeyHash: keyHash,
		Name:    body.Name,
	})
	if err != nil {
		return nil, exception.ErrInternal("failed to create API key")
	}

	return map[string]any{
		"apikey": apikey,
		"key":    rawKey,
	}, nil
}

func (s *AuthService) DeleteAPIKey(ctx *fiber.Ctx) (map[string]any, error) {
	keyID := ctx.Params("id")
	if keyID == "" {
		return nil, exception.ErrBadRequest("missing key id")
	}

	err := s.repo.DeleteAPIKey(context.Background(), keyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNotFound("api key not found")
		}
		return nil, exception.ErrInternal("failed to delete api key")
	}
	return map[string]any{"deleted": keyID}, nil
}

func (s *AuthService) ListAPIKeys(ctx *fiber.Ctx) ([]*dbs.Apikey, error) {
	userID := ctx.Params("user_id")
	if userID == "" {
		return nil, exception.ErrBadRequest("missing user id")
	}

	keys, err := s.repo.GetAPIKeysByUser(context.Background(), userID)
	if err != nil {
		return nil, exception.ErrInternal("failed to fetch API keys")
	}
	return keys, nil
}

// ============ JWT Helpers ============
func (s *AuthService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.jwtTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) extractUserID(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", exception.ErrUnauthorized("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", exception.ErrUnauthorized("invalid Authorization format")
	}

	tokenStr := parts[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", exception.ErrUnauthorized("invalid or expired token")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", exception.ErrUnauthorized("invalid token claims")
	}

	return userID, nil
}

// ============ Password / API Key Helpers ============
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hash)
}

func GenerateAPIKey(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate API key: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func HashAPIKey(key string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash API key: %v", err)
	}
	return string(hash)
}
