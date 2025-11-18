package transfer

// AuthRegisterRequest
type AuthRegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

// AuthLoginRequest
type AuthLoginRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=64"`
	Password  string `json:"password" validate:"required,min=6,max=72"`
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// AuthRefreshRequest
type AuthRefreshRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthLogoutRequest
type AuthLogoutRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

// ======== AUTH RESPONSES ========

// AuthUserResponse
type AuthUserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	IsBlocked bool   `json:"is_blocked"`
	IsChecked bool   `json:"is_checked"`
}

// AuthSessionResponse
type AuthSessionResponse struct {
	ID            string `json:"id"`
	AccessJti     string `json:"access_jti"`
	RefreshJti    string `json:"refresh_jti"`
	AccessExp     string `json:"access_exp"`
	RefreshExp    string `json:"refresh_exp"`
	AccessStatus  string `json:"access_status"`
	RefreshStatus string `json:"refresh_status"`
	IpAddress     string `json:"ip_address"`
	UserAgent     string `json:"user_agent"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// AuthLoginResponse
type AuthLoginResponse struct {
	User         *AuthUserResponse    `json:"user"`
	Session      *AuthSessionResponse `json:"session"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	TokenType    string               `json:"token_type" default:"Bearer"`
	ExpiresIn    int64                `json:"expires_in"`
}

// AuthMeResponse
type AuthMeResponse struct {
	User     *AuthUserResponse      `json:"user"`
	Sessions []*AuthSessionResponse `json:"sessions"`
}

// AuthTokenResponse
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type" default:"Bearer"`
	ExpiresIn    int64  `json:"expires_in"`
}

// AuthClaims
type AuthClaims struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	Session    string `json:"session"`
	AccessJti  string `json:"access_jti"`
	RefreshJti string `json:"refresh_jti"`
}
