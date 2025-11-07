package mapper

import (
    "brickbang/storage/dbs"
)

// DTOs
type RoleCreateRequest struct {
    // TODO: заполнить реальные поля
}

type RoleUpdateRequest struct {
    // TODO: заполнить реальные поля
}

type RoleResponse struct {
    ID string `json:"id"`
    // TODO: добавить остальные поля
}

// Mappers
func RoleToCreateParams(req *RoleCreateRequest) *dbs. {
    return &dbs.{
	// TODO: map fields
    }
}

func RoleToUpdateParams(id string, req *RoleUpdateRequest) *dbs.UpdateRoleByIDParams {
    return &dbs.UpdateRoleByIDParams{
	ID: id,
	// TODO: map fields
    }
}

func RoleToResponse(row *dbs.) *RoleResponse {
    return &RoleResponse{
	ID: row.ID,
	// TODO: map other fields
    }
}

func RoleToResponseList(rows []*dbs.) []*RoleResponse {
    list := make([]*RoleResponse, 0, len(rows))
    for _, r := range rows {
	list = append(list, RoleToResponse(r))
    }
    return list
}
