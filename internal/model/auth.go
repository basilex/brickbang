package model

// AuthLoginRequest — вход по логину и паролю.
type AuthLoginRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=64"`
	Password  string `json:"password" validate:"required,min=6,max=128"`
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// AuthRegisterRequest — регистрация нового пользователя.
type AuthRegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

// AuthRefreshRequest — обновление пары токенов.
type AuthRefreshRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthLogoutRequest — завершение активной сессии.
type AuthLogoutRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

// ======== AUTH RESPONSES ========

// AuthUserResponse — минимальные данные о пользователе.
type AuthUserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	IsBlocked bool   `json:"is_blocked"`
	IsChecked bool   `json:"is_checked"`
}

// AuthSessionResponse — сведения об активной сессии пользователя.
type AuthSessionResponse struct {
	ID            string `json:"id"`
	IpAddress     string `json:"ip_address"`
	UserAgent     string `json:"user_agent"`
	AccessExp     string `json:"access_exp"`
	RefreshExp    string `json:"refresh_exp"`
	AccessStatus  string `json:"access_status"`
	RefreshStatus string `json:"refresh_status"`
	CreatedAt     string `json:"created_at"`
}

// AuthLoginResponse — ответ логина.
type AuthLoginResponse struct {
	User    *AuthUserResponse    `json:"user"`
	Session *AuthSessionResponse `json:"session"`
}

// AuthMeResponse — объединённый ответ при запросе /me.
type AuthMeResponse struct {
	User     *AuthUserResponse      `json:"user"`
	Sessions []*AuthSessionResponse `json:"sessions"`
}

// AuthTokenResponse — ответ с новой парой токенов.
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type" default:"Bearer"`
	ExpiresIn    int64  `json:"expires_in"`
}

// AuthClaims — структура для JWT claims (внутренняя).
type AuthClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Session  string `json:"session"`
}
