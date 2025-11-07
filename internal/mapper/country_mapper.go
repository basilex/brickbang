package mapper

import (
    "brickbang/storage/dbs"
)

// DTOs
type CountryCreateRequest struct {
    // TODO: заполнить реальные поля
}

type CountryUpdateRequest struct {
    // TODO: заполнить реальные поля
}

type CountryResponse struct {
    ID string `json:"id"`
    // TODO: добавить остальные поля
}

// Mappers
func CountryToCreateParams(req *CountryCreateRequest) *dbs.CreateCountryParams {
    return &dbs.CreateCountryParams{
	// TODO: map fields
    }
}

func CountryToUpdateParams(id string, req *CountryUpdateRequest) *dbs.UpdateCountryByIDParams {
    return &dbs.UpdateCountryByIDParams{
	ID: id,
	// TODO: map fields
    }
}

func CountryToResponse(row *dbs.) *CountryResponse {
    return &CountryResponse{
	ID: row.ID,
	// TODO: map other fields
    }
}

func CountryToResponseList(rows []*dbs.) []*CountryResponse {
    list := make([]*CountryResponse, 0, len(rows))
    for _, r := range rows {
	list = append(list, CountryToResponse(r))
    }
    return list
}
