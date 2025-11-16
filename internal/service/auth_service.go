package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"brickbang/internal/config"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/transfer"
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
	Register(ctx context.Context, req *transfer.AuthRegisterRequest) (*transfer.AuthUserResponse, error)
	Login(ctx context.Context, req *transfer.AuthLoginRequest) (*transfer.AuthLoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, userID, refreshToken string) (*transfer.AuthLoginResponse, error)
	Me(ctx context.Context, userID string) (*transfer.AuthMeResponse, error)
	Block(ctx context.Context, userID string, blocked bool) (*transfer.AuthUserResponse, error)
}

type authService struct {
	queries   *dbs.Queries
	blacklist IBlacklistService
}

func NewAuthService(q *dbs.Queries, bl IBlacklistService) IAuthService {
	return &authService{queries: q, blacklist: bl}
}

// --- Блокировка access токена по JTI
func (s *authService) blockAccessTokenByJTI(ctx context.Context, accessJTI string, accessExp time.Time) error {
	if accessJTI == "" {
		return nil
	}
	if time.Now().UTC().After(accessExp) {
		_, err := s.queries.ExpireAccessTokenByJTI(ctx, accessJTI)
		return err
	}
	_, err := s.queries.RevokeAccessTokenByJTI(ctx, accessJTI)
	return err
}

// --- Register
func (s *authService) Register(ctx context.Context, req *transfer.AuthRegisterRequest) (*transfer.AuthUserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.CreateUser(ctx, &dbs.CreateUserParams{
		Username:  req.Username,
		Password:  string(hashed),
		IsChecked: false,
	})
	if err != nil {
		return nil, err
	}

	return &transfer.AuthUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsBlocked: user.IsBlocked,
		IsChecked: user.IsChecked,
	}, nil
}

// --- Login
func (s *authService) Login(ctx context.Context, req *transfer.AuthLoginRequest) (*transfer.AuthLoginResponse, error) {
	user, err := s.queries.GetUserByUsername(ctx, req.Username)
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

	// --- Ревокация старых сессий
	sessions, _ := s.queries.ListSessionsByUserID(ctx, user.ID)
	for _, sess := range sessions {
		_ = s.blockAccessTokenByJTI(ctx, sess.AccessJti, sess.AccessExp.Time)
		_, _ = s.queries.RevokeSessionByID(ctx, sess.ID)
	}

	now := time.Now().UTC()
	accessExp := now.Add(config.Get().JWTAccessExpiration)
	refreshExp := now.Add(config.Get().JWTRefreshExpiration)

	// --- Генерация refresh токена с xid
	refreshPlain, _ := utility.GenerateRandomString(32)
	refreshJti := xid.New().String()
	refreshHash := sha256.Sum256([]byte(refreshPlain))
	refreshHashHex := hex.EncodeToString(refreshHash[:])

	// --- Создание сессии
	session, err := s.queries.CreateSession(ctx, &dbs.CreateSessionParams{
		UserID:       user.ID,
		AccessToken:  "",
		AccessJti:    "",
		RefreshToken: refreshHashHex,
		RefreshJti:   refreshJti,
		AccessExp:    utility.ToPGTimestamp(accessExp),
		RefreshExp:   utility.ToPGTimestamp(refreshExp),
		IpAddress:    req.IpAddress,
		UserAgent:    req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// --- Генерация access токена
	accessToken, accessJti, _ := utility.GenerateAccessToken(user.ID, session.ID, config.Get().JWTAccessExpiration)

	updatedSession, err := s.queries.UpdateAccessTokenByID(ctx, &dbs.UpdateAccessTokenByIDParams{
		ID:          session.ID,
		AccessToken: accessToken,
		AccessJti:   accessJti,
		AccessExp:   utility.ToPGTimestamp(accessExp),
	})
	if err != nil {
		_, _ = s.queries.RevokeSessionByID(ctx, session.ID)
		return nil, err
	}

	resp := &transfer.AuthLoginResponse{
		User: &transfer.AuthUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			IsBlocked: user.IsBlocked,
			IsChecked: user.IsChecked,
		},
		Session: &transfer.AuthSessionResponse{
			ID:            updatedSession.ID,
			AccessJti:     updatedSession.AccessJti,
			RefreshJti:    updatedSession.RefreshJti,
			AccessExp:     utility.FromPGTimestampToString(updatedSession.AccessExp),
			RefreshExp:    utility.FromPGTimestampToString(updatedSession.RefreshExp),
			AccessStatus:  updatedSession.AccessStatus,
			RefreshStatus: updatedSession.RefreshStatus,
			IpAddress:     updatedSession.IpAddress,
			UserAgent:     updatedSession.UserAgent,
			CreatedAt:     utility.FromPGTimestampToString(updatedSession.CreatedAt),
			UpdatedAt:     utility.FromPGTimestampToString(updatedSession.UpdatedAt),
		},
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
		ExpiresIn:    int64(config.Get().JWTAccessExpiration.Seconds()),
	}

	return resp, nil
}

// --- Logout
func (s *authService) Logout(ctx context.Context, sessionID string) error {
	sess, err := s.queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	return s.blockAccessTokenByJTI(ctx, sess.AccessJti, sess.AccessExp.Time)
}

// --- Refresh
func (s *authService) Refresh(ctx context.Context, userID, refreshToken string) (*transfer.AuthLoginResponse, error) {
	sessions, err := s.queries.ListSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(refreshToken))
	refreshHash := hex.EncodeToString(hash[:])

	var sess *dbs.Session
	for i := range sessions {
		if sessions[i].RefreshStatus != "valid" {
			continue
		}
		if sessions[i].RefreshToken == refreshHash {
			sess = sessions[i]
			break
		}
	}

	if sess == nil {
		return nil, ErrInvalidRefreshToken
	}

	if !sess.RefreshExp.Valid || sess.RefreshExp.Time.Before(time.Now().UTC()) {
		return nil, ErrRefreshTokenExpired
	}

	_ = s.blockAccessTokenByJTI(ctx, sess.AccessJti, sess.AccessExp.Time)

	// --- Генерация новых токенов с xid
	newAccessToken, newAccessJti, _ := utility.GenerateAccessToken(sess.UserID, sess.ID, config.Get().JWTAccessExpiration)
	newRefreshPlain, _ := utility.GenerateRandomString(32)
	newRefreshJti := xid.New().String()
	newRefreshHash := sha256.Sum256([]byte(newRefreshPlain))
	newRefreshHashHex := hex.EncodeToString(newRefreshHash[:])

	newAccessExp := time.Now().UTC().Add(config.Get().JWTAccessExpiration)
	newRefreshExp := time.Now().UTC().Add(config.Get().JWTRefreshExpiration)

	updatedSession, err := s.queries.RotateTokensByID(ctx, &dbs.RotateTokensByIDParams{
		ID:           sess.ID,
		AccessToken:  newAccessToken,
		AccessJti:    newAccessJti,
		AccessExp:    pgtype.Timestamp{Time: newAccessExp, Valid: true},
		RefreshToken: newRefreshHashHex,
		RefreshJti:   newRefreshJti,
		RefreshExp:   pgtype.Timestamp{Time: newRefreshExp, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	formatTime := func(ts pgtype.Timestamp) string {
		if ts.Valid {
			return ts.Time.UTC().Format(time.RFC3339)
		}
		return ""
	}

	sessionResp := &transfer.AuthSessionResponse{
		ID:            updatedSession.ID,
		AccessJti:     updatedSession.AccessJti,
		AccessExp:     formatTime(updatedSession.AccessExp),
		RefreshJti:    updatedSession.RefreshJti,
		RefreshExp:    formatTime(updatedSession.RefreshExp),
		AccessStatus:  updatedSession.AccessStatus,
		RefreshStatus: updatedSession.RefreshStatus,
		IpAddress:     updatedSession.IpAddress,
		UserAgent:     updatedSession.UserAgent,
		CreatedAt:     formatTime(updatedSession.CreatedAt),
		UpdatedAt:     formatTime(updatedSession.UpdatedAt),
	}

	return &transfer.AuthLoginResponse{
		User: &transfer.AuthUserResponse{
			ID: sess.UserID,
		},
		Session:      sessionResp,
		TokenType:    "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshPlain,
		ExpiresIn:    int64(config.Get().JWTAccessExpiration.Seconds()),
	}, nil
}

// --- Me
func (s *authService) Me(ctx context.Context, userID string) (*transfer.AuthMeResponse, error) {
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessionsDb, _ := s.queries.ListSessionsByUserID(ctx, userID)
	sessions := make([]*transfer.AuthSessionResponse, 0, len(sessionsDb))

	for _, sess := range sessionsDb {
		sessions = append(sessions, &transfer.AuthSessionResponse{
			ID:            sess.ID,
			AccessJti:     sess.AccessJti,
			AccessExp:     utility.FromPGTimestampToString(sess.AccessExp),
			RefreshJti:    sess.RefreshJti,
			RefreshExp:    utility.FromPGTimestampToString(sess.RefreshExp),
			AccessStatus:  sess.AccessStatus,
			RefreshStatus: sess.RefreshStatus,
			IpAddress:     sess.IpAddress,
			UserAgent:     sess.UserAgent,
			CreatedAt:     utility.FromPGTimestampToString(sess.CreatedAt),
			UpdatedAt:     utility.FromPGTimestampToString(sess.UpdatedAt),
		})
	}

	userResp := &transfer.AuthUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsBlocked: user.IsBlocked,
		IsChecked: user.IsChecked,
	}

	return &transfer.AuthMeResponse{
		User:     userResp,
		Sessions: sessions,
	}, nil
}

// --- Block / Unblock
func (s *authService) Block(ctx context.Context, userID string, blocked bool) (*transfer.AuthUserResponse, error) {
	user, err := s.queries.UpdateUserIsBlockedByID(ctx, &dbs.UpdateUserIsBlockedByIDParams{
		ID:        userID,
		IsBlocked: blocked,
	})
	if err != nil {
		return nil, err
	}

	if blocked {
		sessions, _ := s.queries.ListSessionsByUserID(ctx, userID)
		for _, sess := range sessions {
			_ = s.blockAccessTokenByJTI(ctx, sess.AccessJti, sess.AccessExp.Time)
			_, _ = s.queries.RevokeSessionByID(ctx, sess.ID)
		}
	}

	return &transfer.AuthUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsBlocked: user.IsBlocked,
		IsChecked: user.IsChecked,
	}, nil
}
