package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"brickbang/internal/config"
	"brickbang/internal/mapper"
	"brickbang/internal/model"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/utility"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserBlocked         = errors.New("user is blocked")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrUserNotFound        = errors.New("user not found")
)

type IAuthService interface {
	Register(ctx context.Context, req *model.AuthRegisterRequest) (*model.AuthUserResponse, error)
	Login(ctx context.Context, req *model.AuthLoginRequest) (*model.AuthLoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, userID, refreshToken string) (*model.AuthSessionResponse, error)
	Me(ctx context.Context, userID string) (*model.AuthMeResponse, error)
	Block(ctx context.Context, userID string, blocked bool) (*model.AuthUserResponse, error)
}

type AuthService struct {
	q *dbs.Queries
}

func NewAuthService(q *dbs.Queries) IAuthService {
	return &AuthService{q: q}
}

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

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	accessExp := now.Add(config.Get().JWTAccessExpiration)
	refreshExp := now.Add(config.Get().JWTRefreshExpiration)

	// Генерация токенов
	accessToken, err := utility.GenerateAccessToken(u.ID, "", time.Hour)
	if err != nil {
		return nil, err
	}

	refreshTokenPlain, err := utility.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshTokenPlain), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создание сессии
	session, err := s.q.AuthCreateSession(ctx, &dbs.AuthCreateSessionParams{
		UserID:       u.ID,
		AccessToken:  accessToken,
		RefreshToken: string(refreshTokenHash),
		AccessExp:    utility.ToPGTimestamp(accessExp),
		RefreshExp:   utility.ToPGTimestamp(refreshExp),
		IpAddress:    req.IpAddress,
		UserAgent:    req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// Ответ клиенту — plaintext refresh токен возвращаем только один раз
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
			RefreshToken:  refreshTokenPlain, // <- оригинал
		}),
	}

	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	_, err := s.q.AuthRevokeSessionByID(ctx, sessionID)
	return err
}

func (s *AuthService) Refresh(ctx context.Context, userID, refreshToken string) (*model.AuthSessionResponse, error) {
	sessions, err := s.q.AuthListSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sess *dbs.Session
	for i := range sessions {
		if sessions[i].RefreshStatus != "valid" {
			continue
		}
		if bcrypt.CompareHashAndPassword([]byte(sessions[i].RefreshToken), []byte(refreshToken)) == nil {
			sess = sessions[i]
			break
		}
	}

	if sess == nil {
		return nil, ErrInvalidRefreshToken
	}

	if sess.RefreshExp.Time.Before(time.Now().UTC()) {
		return nil, ErrRefreshTokenExpired
	}

	newAccess, err := utility.GenerateAccessToken(sess.UserID, sess.ID, time.Hour)
	if err != nil {
		return nil, err
	}
	newAccessExp := time.Now().UTC().Add(time.Hour)

	updated, err := s.q.AuthUpdateAccessTokenByID(ctx, &dbs.AuthUpdateAccessTokenByIDParams{
		ID:          sess.ID,
		AccessToken: newAccess,
		AccessExp:   utility.ToPGTimestamp(newAccessExp),
	})
	if err != nil {
		return nil, err
	}

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
