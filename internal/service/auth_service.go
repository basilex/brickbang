package service

import (
	"brickbang/internal/mapper"
	"brickbang/internal/model"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/utility"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ======== ERRORS ========
var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserBlocked         = errors.New("user is blocked")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrUserNotFound        = errors.New("user not found")
)

// IAuthService — интерфейс сервиса аутентификации
type IAuthService interface {
	Register(ctx context.Context, req *model.AuthRegisterRequest) (*model.AuthUserResponse, error)
	Login(ctx context.Context, req *model.AuthLoginRequest) (*model.AuthLoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, token string) (*model.AuthSessionResponse, error)
	Me(ctx context.Context, userID string) (*model.AuthMeResponse, error)
	Block(ctx context.Context, userID string, blocked bool) (*model.AuthUserResponse, error)
}

// AuthService — реализация сервиса
type AuthService struct {
	q *dbs.Queries
}

// NewAuthService — конструктор
func NewAuthService(q *dbs.Queries) IAuthService {
	return &AuthService{q: q}
}

// Register — регистрация нового пользователя.
func (s *AuthService) Register(ctx context.Context, req *model.AuthRegisterRequest) (*model.AuthUserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.q.AuthCreateUser(ctx, &dbs.AuthCreateUserParams{
		Username:  req.Username,
		Password:  string(hashed),
		IsChecked: true,
	})
	if err != nil {
		return nil, err
	}

	u := &dbs.User{
		ID:        user.ID,
		Username:  user.Username,
		IsBlocked: false,
		IsChecked: user.IsChecked,
	}

	return mapper.MapUserToAuthUserResponse(u), nil
}

func (s *AuthService) Login(ctx context.Context, req *model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	// 1. Получаем пользователя
	u, err := s.q.AuthSelectUserCredentials(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if u.IsBlocked {
		return nil, ErrUserBlocked
	}

	// 2. Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	accessExp := now.Add(time.Hour)
	refreshExp := now.Add(7 * 24 * time.Hour)

	// 3. Генерируем токены
	accessToken, err := utility.GenerateAccessToken(u.ID, "", time.Hour)
	if err != nil {
		return nil, err
	}
	refreshToken, _ := utility.GenerateRandomString(32)

	// 4. Создаем сессию в БД
	session, err := s.q.AuthCreateSession(ctx, &dbs.AuthCreateSessionParams{
		UserID:       u.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExp:    utility.ToPGTimestamp(accessExp),
		RefreshExp:   utility.ToPGTimestamp(refreshExp),
		IpAddress:    req.IpAddress,
		UserAgent:    req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// 5. Формируем DTO через mapper
	resp := &model.AuthLoginResponse{
		User: mapper.MapUserToAuthUserResponse(&dbs.User{
			ID:        u.ID,
			Username:  u.Username,
			IsBlocked: u.IsBlocked,
			IsChecked: u.IsChecked,
		}),
		Session: mapper.MapSessionToAuthSessionResponse(&dbs.Session{
			ID:            session.ID,
			IpAddress:     session.IpAddress,
			UserAgent:     session.UserAgent,
			AccessExp:     session.AccessExp,
			RefreshExp:    session.RefreshExp,
			AccessStatus:  session.AccessStatus,
			RefreshStatus: session.RefreshStatus,
			CreatedAt:     session.CreatedAt,
			AccessToken:   session.AccessToken,
			RefreshToken:  session.RefreshToken,
		}),
	}

	return resp, nil
}

// Logout — завершение активной сессии.
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	_, err := s.q.AuthRevokeAccessSessionByID(ctx, sessionID)
	return err
}

// Refresh — обновление access токена по refresh токену.
func (s *AuthService) Refresh(ctx context.Context, token string) (*model.AuthSessionResponse, error) {
	// Получаем сессию по refresh токену
	sess, err := s.q.AuthSelectSessionByRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	// Проверяем статус и срок действия refresh токена
	if sess.RefreshStatus != "valid" || sess.RefreshExp.Time.Before(time.Now().UTC()) {
		return nil, ErrRefreshTokenExpired
	}

	// Генерируем новый access токен и новый срок действия
	newAccess, _ := utility.GenerateRandomString(32)
	newAccessExp := time.Now().UTC().Add(time.Hour)

	// Обновляем access токен в БД
	updated, err := s.q.AuthUpdateAccessTokenByID(ctx, &dbs.AuthUpdateAccessTokenByIDParams{
		ID:          sess.ID,
		AccessToken: newAccess,
		AccessExp:   utility.ToPGTimestamp(newAccessExp),
	})
	if err != nil {
		return nil, err
	}

	// Формируем DTO сразу из обновленных и существующих данных сессии
	resp := &model.AuthSessionResponse{
		ID:            updated.ID,
		IpAddress:     sess.IpAddress,
		UserAgent:     sess.UserAgent,
		AccessExp:     updated.AccessExp.Time.Format(time.RFC3339),
		RefreshExp:    sess.RefreshExp.Time.Format(time.RFC3339),
		AccessStatus:  updated.AccessStatus,
		RefreshStatus: sess.RefreshStatus,
		CreatedAt:     sess.CreatedAt.Time.Format(time.RFC3339),
	}

	return resp, nil
}

// Me — возврат информации о пользователе и активных сессиях.
func (s *AuthService) Me(ctx context.Context, userID string) (*model.AuthMeResponse, error) {
	u, err := s.q.AuthSelectUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	sessions, err := s.q.AuthListSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapper.MapUserAndSessionsToMeResponse(&dbs.User{
		ID:        u.ID,
		Username:  u.Username,
		IsBlocked: u.IsBlocked,
		IsChecked: u.IsChecked,
	}, sessions), nil
}

// Block — блокировка или разблокировка пользователя администратором.
func (s *AuthService) Block(ctx context.Context, userID string, blocked bool) (*model.AuthUserResponse, error) {
	u, err := s.q.UserUpdateIsBlockedByID(ctx, &dbs.UserUpdateIsBlockedByIDParams{
		ID:        userID,
		IsBlocked: blocked,
	})
	if err != nil {
		return nil, err
	}

	return mapper.MapUserToAuthUserResponse(&dbs.User{
		ID:        u.ID,
		Username:  u.Username,
		IsBlocked: u.IsBlocked,
		IsChecked: u.IsChecked,
	}), nil
}
