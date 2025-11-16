package model

// RoleCreateRequest
type RoleCreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

// RoleUpdateRequest
type RoleUpdateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

// RoleResponse
type RoleResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
