package service

import (
	"context"
	"time"

	"brickbang/internal/repository"
	"brickbang/storage/dbs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// IAuthService defines methods for authentication and API key management
type IAuthService interface {
	Login(ctx *fiber.Ctx) (map[string]interface{}, error)
	Refresh(ctx *fiber.Ctx) (map[string]interface{}, error)
	Me(ctx *fiber.Ctx) (map[string]interface{}, error)
	CreateAPIKey(ctx *fiber.Ctx) (map[string]interface{}, error)
	DeleteAPIKey(ctx *fiber.Ctx) (map[string]interface{}, error)
	ListAPIKeys(ctx *fiber.Ctx) ([]*dbs.Apikey, error)
}

// ===== Implementation ======
type AuthService struct {
	repo      repository.IAuthRepository
	jwtSecret string
	jwtTTL    time.Duration
}

// NewAuthService constructor
func NewAuthService(repo repository.IAuthRepository, jwtSecret string, jwtTTL time.Duration) IAuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

// ===== JWT login ======
func (s *AuthService) Login(ctx *fiber.Ctx) (map[string]interface{}, error) {
	type loginBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body loginBody
	if err := ctx.BodyParser(&body); err != nil {
		return nil, fiber.ErrBadRequest
	}

	user, err := s.repo.FindByUsername(context.Background(), body.Username)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	// Проверка пароля (bcrypt)
	if !CheckPasswordHash(body.Password, user.Password) {
		return nil, fiber.ErrUnauthorized
	}

	// Обновляем visited_at
	_, _ = s.repo.UpdateVisitedAt(context.Background(), user.ID)

	// Генерируем JWT
	tokenString, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
		"token": tokenString,
	}, nil
}

// ===== JWT refresh ======
func (s *AuthService) Refresh(ctx *fiber.Ctx) (map[string]interface{}, error) {
	// TBD: парсим Bearer, проверяем валидность и выдаем новый токен
	return map[string]interface{}{
		"message": "refresh not implemented yet",
	}, nil
}

// ===== Current user ======
func (s *AuthService) Me(ctx *fiber.Ctx) (map[string]interface{}, error) {
	// TBD: берем userID из JWT
	return map[string]interface{}{
		"message": "me not implemented yet",
	}, nil
}

// ===== API Key management ======
func (s *AuthService) CreateAPIKey(ctx *fiber.Ctx) (map[string]interface{}, error) {
	type req struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}

	var body req
	if err := ctx.BodyParser(&body); err != nil {
		return nil, fiber.ErrBadRequest
	}

	keyHash := GenerateAPIKeyHash() // TODO: генерация случайного API key

	apikey, err := s.repo.CreateAPIKey(context.Background(), &dbs.AuthCreateAPIKeyParams{
		UserID:  body.UserID,
		KeyHash: keyHash,
		Name:    body.Name,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"apikey": apikey,
		"key":    keyHash, // показываем пользователю только один раз
	}, nil
}

func (s *AuthService) DeleteAPIKey(ctx *fiber.Ctx) (map[string]interface{}, error) {
	keyID := ctx.Params("id")
	if keyID == "" {
		return nil, fiber.ErrBadRequest
	}

	if err := s.repo.DeleteAPIKey(context.Background(), keyID); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"deleted": keyID,
	}, nil
}

func (s *AuthService) ListAPIKeys(ctx *fiber.Ctx) ([]*dbs.Apikey, error) {
	userID := ctx.Params("user_id")
	if userID == "" {
		return nil, fiber.ErrBadRequest
	}

	return s.repo.GetAPIKeysByUser(context.Background(), userID)
}

// ===== JWT helper ======
func (s *AuthService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.jwtTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ===== Password helper ======
func CheckPasswordHash(password, hash string) bool {
	// TODO: использовать bcrypt.CompareHashAndPassword
	return password == hash
}

// ===== API Key helper ======
func GenerateAPIKeyHash() string {
	// TODO: сгенерировать безопасный случайный ключ, вернуть как строку
	return "TODO-API-KEY"
}
