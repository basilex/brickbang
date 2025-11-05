package mapper

import (
    "time"

    "brickbang/storage/dbs"
)

type {{.Entity}}CreateRequest struct {
    Name string `json:"name" validate:"required,min=2,max=100"`
    Code string `json:"code" validate:"required,len=3,uppercase"`
}

type {{.Entity}}UpdateRequest struct {
    Name string `json:"name" validate:"omitempty,min=2,max=100"`
    Code string `json:"code" validate:"omitempty,len=3,uppercase"`
}

type {{.Entity}}Response struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Code      string    `json:"code"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func To{{.Entity}}NewParams(req *{{.Entity}}CreateRequest) *dbs.{{.Entity}}NewParams {
    return &dbs.{{.Entity}}NewParams{
	Name:      req.Name,
	Code:      req.Code,
	CreatedAt: time.Now(),
    }
}

func To{{.Entity}}UpdateParams(id string, req *{{.Entity}}UpdateRequest) *dbs.{{.Entity}}UpdateByIDParams {
    params := &dbs.{{.Entity}}UpdateByIDParams{
	    ID:        id,
	    Name:      req.Name,
	    Code:      req.Code,
	    UpdatedAt: time.Now(),
    }
    return params
}

func To{{.Entity}}Response(row *dbs.{{.Entity}}) *{{.Entity}}Response {
    return &{{.Entity}}Response{
	    ID:        row.ID,
	    Name:      row.Name,
	    Code:      row.Code,
	    CreatedAt: row.CreatedAt,
	    UpdatedAt: row.UpdatedAt,
    }
}

func To{{.Entity}}ResponseList(rows []dbs.{{.Entity}}) []*{{.Entity}}Response {
    if len(rows) == 0 {
	    return []*{{.Entity}}Response{}
    }

    result := make([]*{{.Entity}}Response, 0, len(rows))
    for _, r := range rows {
	    result = append(result, To{{.Entity}}Response(&r))
    }
    return result
}
