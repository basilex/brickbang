package transfer

// GrantCreateRequest
type GrantCreateRequest struct {
	Code        string `json:"code" validate:"required,min=5,max=255"`
	Description string `json:"description" validate:"max=1024"`
}

// GrantUpdateRequest
type GrantUpdateRequest struct {
	Code        string `json:"code" validate:"required,min=5,max=255"`
	Description string `json:"description" validate:"max=1024"`
}

// GrantResponse
type GrantResponse struct {
	Pid         string `json:"pid"`
	Code        string `json:"code"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
