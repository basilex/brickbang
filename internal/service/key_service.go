package service

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/exception"
	"brickbang/internal/repository"
	"brickbang/storage/dbs"
)

type IKeyService interface {
	CreateAPIKey(ctx *fiber.Ctx) (map[string]any, error)
	DeleteAPIKey(ctx *fiber.Ctx) (map[string]any, error)
	ListAPIKeys(ctx *fiber.Ctx) ([]*dbs.Apikey, error)
}

type KeyService struct {
	repo repository.IKeyRepository
}

func NewKeyService(repo repository.IKeyRepository) IKeyService {
	return &KeyService{repo: repo}
}

func (rcv *KeyService) CreateAPIKey(ctx *fiber.Ctx) (map[string]any, error) {
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

	apikey, err := rcv.repo.CreateAPIKey(context.Background(), &dbs.AuthCreateAPIKeyParams{
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

func (rcv *KeyService) DeleteAPIKey(ctx *fiber.Ctx) (map[string]any, error) {
	keyID := ctx.Params("id")
	if keyID == "" {
		return nil, exception.ErrBadRequest("missing key id")
	}

	err := rcv.repo.DeleteAPIKey(context.Background(), keyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNotFound("api key not found")
		}
		return nil, exception.ErrInternal("failed to delete api key")
	}
	return map[string]any{"deleted": keyID}, nil
}

func (rcv *KeyService) ListAPIKeys(ctx *fiber.Ctx) ([]*dbs.Apikey, error) {
	userID := ctx.Params("user_id")
	if userID == "" {
		return nil, exception.ErrBadRequest("missing user id")
	}

	keys, err := rcv.repo.APIKeysByUser(context.Background(), userID)
	if err != nil {
		return nil, exception.ErrInternal("failed to fetch API keys")
	}
	return keys, nil
}
