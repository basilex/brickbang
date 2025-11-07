package mapper

import (
    "{{.ModulePath}}/storage/dbs"
)

// DTOs
type {{.Entity}}CreateRequest struct {
    // TODO: заполнить реальные поля
}

type {{.Entity}}UpdateRequest struct {
    // TODO: заполнить реальные поля
}

type {{.Entity}}Response struct {
    ID string `json:"id"`
    // TODO: добавить остальные поля
}

// Mappers
func {{.Entity}}ToCreateParams(req *{{.Entity}}CreateRequest) *dbs.{{.CreateParams}} {
    return &dbs.{{.CreateParams}}{
	// TODO: map fields
    }
}

func {{.Entity}}ToUpdateParams(id string, req *{{.Entity}}UpdateRequest) *dbs.{{.UpdateParams}} {
    return &dbs.{{.UpdateParams}}{
	ID: id,
	// TODO: map fields
    }
}

func {{.Entity}}ToResponse(row *dbs.{{.ModelStruct}}) *{{.Entity}}Response {
    return &{{.Entity}}Response{
	ID: row.ID,
	// TODO: map other fields
    }
}

func {{.Entity}}ToResponseList(rows []*dbs.{{.ModelStruct}}) []*{{.Entity}}Response {
    list := make([]*{{.Entity}}Response, 0, len(rows))
    for _, r := range rows {
	list = append(list, {{.Entity}}ToResponse(r))
    }
    return list
}
