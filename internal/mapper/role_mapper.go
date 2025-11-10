package mapper

import (
	"time"

	"brickbang/internal/repository/dbs"
)

// DTOs
type RoleCreateRequest struct {
	Name string
}

type RoleUpdateRequest struct {
	Name string
}

type RoleResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Mappers
func RoleToCreateParams(req *RoleCreateRequest) *dbs.Role {
	return &dbs.Role{
		Name: req.Name,
	}
}

func RoleToUpdateParams(id string, req *RoleUpdateRequest) *dbs.UpdateRoleByIDParams {
	return &dbs.UpdateRoleByIDParams{
		ID:   id,
		Name: req.Name,
	}
}

func RoleToResponse(row *dbs.Role) *RoleResponse {
	return &RoleResponse{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func RoleToResponseList(rows []*dbs.Role) []*RoleResponse {
	list := make([]*RoleResponse, 0, len(rows))
	for _, r := range rows {
		list = append(list, RoleToResponse(r))
	}
	return list
}
