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

type authService struct {
	queries   *dbs.Queries
	blacklist IBlacklistService
}

func NewAuthService(q *dbs.Queries, bl IBlacklistService) IAuthService {
	return &authService{queries: q, blacklist: bl}
}

// --- private helper
func (rcv *authService) blockAccessToken(ctx context.Context, token string, exp time.Time) {
	if token == "" {
		return
	}
	if claims, err := utility.ParseAccessToken(token); err == nil {
		ttl := max(time.Until(exp), 0)
		_ = rcv.blacklist.BlockToken(ctx, claims.JTI, ttl)
	}
}

// --- Register
func (rcv *authService) Register(ctx context.Context, req *model.AuthRegisterRequest) (*model.AuthUserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := rcv.queries.AuthCreateUser(ctx, &dbs.AuthCreateUserParams{
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

// --- Login
func (rcv *authService) Login(ctx context.Context, req *model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	user, err := rcv.queries.AuthSelectUserCredentials(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if user.IsBlocked {
		return nil, ErrUserBlocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	accessExp := now.Add(config.Get().JWTAccessExpiration)
	refreshExp := now.Add(config.Get().JWTRefreshExpiration)

	randStr, err := utility.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	refreshPlain, _, err := utility.GenerateRefreshToken(randStr)
	if err != nil {
		return nil, err
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshPlain), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	session, err := rcv.queries.AuthCreateSession(ctx, &dbs.AuthCreateSessionParams{
		UserID:       user.ID,
		AccessToken:  "",
		RefreshToken: string(refreshTokenHash),
		AccessExp:    utility.ToPGTimestamp(accessExp),
		RefreshExp:   utility.ToPGTimestamp(refreshExp),
		IpAddress:    req.IpAddress,
		UserAgent:    req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	accessToken, _, err := utility.GenerateAccessToken(user.ID, session.ID, config.Get().JWTAccessExpiration)
	if err != nil {
		_, _ = rcv.queries.AuthRevokeSessionByID(ctx, session.ID)
		return nil, err
	}

	updated, err := rcv.queries.AuthUpdateAccessTokenByID(ctx, &dbs.AuthUpdateAccessTokenByIDParams{
		ID:          session.ID,
		AccessToken: accessToken,
		AccessExp:   utility.ToPGTimestamp(accessExp),
	})
	if err != nil {
		_, _ = rcv.queries.AuthRevokeSessionByID(ctx, session.ID)
		return nil, err
	}

	resp := &model.AuthLoginResponse{
		User: mapper.MapUserToAuthUserResponse(&dbs.User{
			ID:        user.ID,
			Username:  user.Username,
			IsBlocked: user.IsBlocked,
			IsChecked: user.IsChecked,
		}),
		Session: mapper.MapSessionToAuthSessionResponse(&dbs.Session{
			ID:            updated.ID,
			IpAddress:     session.IpAddress,
			UserAgent:     session.UserAgent,
			AccessExp:     updated.AccessExp,
			RefreshExp:    session.RefreshExp,
			AccessStatus:  updated.AccessStatus,
			RefreshStatus: session.RefreshStatus,
			CreatedAt:     session.CreatedAt,
			AccessToken:   accessToken,
			RefreshToken:  refreshPlain,
		}),
	}

	// _ = accessJTI
	return resp, nil
}

// --- Logout
func (rcv *authService) Logout(ctx context.Context, sessionID string) error {
	sess, err := rcv.queries.AuthSelectSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	rcv.blockAccessToken(ctx, sess.AccessToken, sess.AccessExp.Time)
	_, err = rcv.queries.AuthRevokeSessionByID(ctx, sessionID)
	return err
}

// --- Refresh
func (rcv *authService) Refresh(ctx context.Context, userID, refreshToken string) (*model.AuthSessionResponse, error) {
	sessions, err := rcv.queries.AuthListSessionsByUserID(ctx, userID)
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

	if refreshJTI, err := utility.ParseRefreshTokenPlain(refreshToken); err == nil {
		blocked, err := rcv.blacklist.IsBlocked(ctx, refreshJTI)
		if err != nil {
			return nil, err
		}
		if blocked {
			return nil, ErrInvalidRefreshToken
		}
	}

	if sess.RefreshExp.Time.Before(time.Now().UTC()) {
		return nil, ErrRefreshTokenExpired
	}

	rcv.blockAccessToken(ctx, sess.AccessToken, sess.AccessExp.Time)

	newAccessToken, newJTI, err := utility.GenerateAccessToken(sess.UserID, sess.ID, config.Get().JWTAccessExpiration)
	if err != nil {
		return nil, err
	}
	_ = newJTI

	newAccessExp := time.Now().UTC().Add(config.Get().JWTAccessExpiration)
	updated, err := rcv.queries.AuthUpdateAccessTokenByID(ctx, &dbs.AuthUpdateAccessTokenByIDParams{
		ID:          sess.ID,
		AccessToken: newAccessToken,
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

// --- Me
func (rcv *authService) Me(ctx context.Context, userID string) (*model.AuthMeResponse, error) {
	user, err := rcv.queries.AuthSelectUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	sessions, err := rcv.queries.AuthListSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapper.MapUserAndSessionsToMeResponse(&dbs.User{
		ID:        user.ID,
		Username:  user.Username,
		IsBlocked: user.IsBlocked,
		IsChecked: user.IsChecked,
	}, sessions), nil
}

// --- Block / Unblock
func (rcv *authService) Block(ctx context.Context, userID string, blocked bool) (*model.AuthUserResponse, error) {
	user, err := rcv.queries.UserUpdateIsBlockedByID(ctx, &dbs.UserUpdateIsBlockedByIDParams{
		ID:        userID,
		IsBlocked: blocked,
	})
	if err != nil {
		return nil, err
	}

	if blocked {
		sessions, err := rcv.queries.AuthListSessionsByUserID(ctx, userID)
		if err == nil {
			for _, s := range sessions {
				rcv.blockAccessToken(ctx, s.AccessToken, s.AccessExp.Time)
				_, _ = rcv.queries.AuthRevokeSessionByID(ctx, s.ID)
			}
		}
	}

	return mapper.MapUserToAuthUserResponse(&dbs.User{
		ID:        user.ID,
		Username:  user.Username,
		IsBlocked: user.IsBlocked,
		IsChecked: user.IsChecked,
	}), nil
}
