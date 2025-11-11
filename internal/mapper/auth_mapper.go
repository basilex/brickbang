package mapper

import (
	"time"

	"brickbang/internal/model"
	"brickbang/internal/repository/dbs"
)

// MapUserToAuthUserResponse — маппинг SQLC-пользователя в DTO.
func MapUserToAuthUserResponse(u *dbs.User) *model.AuthUserResponse {
	return &model.AuthUserResponse{
		ID:        u.ID,
		Username:  u.Username,
		IsBlocked: u.IsBlocked,
		IsChecked: u.IsChecked,
	}
}

// MapSessionToAuthSessionResponse — маппинг SQLC-сессии в DTO.
func MapSessionToAuthSessionResponse(s *dbs.Session) *model.AuthSessionResponse {
	return &model.AuthSessionResponse{
		ID:            s.ID,
		IpAddress:     s.IpAddress,
		UserAgent:     s.UserAgent,
		AccessExp:     formatTimestamp(s.AccessExp.Time),
		RefreshExp:    formatTimestamp(s.RefreshExp.Time),
		AccessStatus:  s.AccessStatus,
		RefreshStatus: s.RefreshStatus,
		CreatedAt:     formatTimestamp(s.CreatedAt.Time),
	}
}

// MapSessionsToAuthSessionResponses — маппинг массива SQLC-сессий в DTO.
func MapSessionsToAuthSessionResponses(sessions []*dbs.Session) []*model.AuthSessionResponse {
	out := make([]*model.AuthSessionResponse, 0, len(sessions))
	for _, s := range sessions {
		out = append(out, MapSessionToAuthSessionResponse(s))
	}
	return out
}

// MapUserAndSessionsToMeResponse — объединяет пользователя и список сессий для ответа /me.
func MapUserAndSessionsToMeResponse(u *dbs.User, sessions []*dbs.Session) *model.AuthMeResponse {
	return &model.AuthMeResponse{
		User:     MapUserToAuthUserResponse(u),
		Sessions: MapSessionsToAuthSessionResponses(sessions),
	}
}

//
// ======================
// Helpers
// ======================
//

// formatTimestamp — безопасное форматирование времени в ISO-8601.
func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
