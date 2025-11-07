package repository

import (
    "context"

    "{{.ModulePath}}/storage/dbs"
)

type I{{.Entity}}Repository interface {
    Create(ctx context.Context, params *dbs.{{.CreateParams}}) (*dbs.{{.ModelStruct}}, error)
    List(ctx context.Context, params *dbs.{{.ListParams}}) ([]*dbs.{{.ModelStruct}}, error)
    GetByID(ctx context.Context, id string) (*dbs.{{.ModelStruct}}, error)
    UpdateByID(ctx context.Context, params *dbs.{{.UpdateParams}}) (*dbs.{{.ModelStruct}}, error)
    DeleteByID(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}

type {{.Entity}}Repository struct {
    q *dbs.Queries
}

func New{{.Entity}}Repository(q *dbs.Queries) I{{.Entity}}Repository {
    return &{{.Entity}}Repository{q: q}
}

func (r *{{.Entity}}Repository) Create(ctx context.Context, params *dbs.{{.CreateParams}}) (*dbs.{{.ModelStruct}}, error) {
    return r.q.Create{{.Entity}}(ctx, *params)
}

func (r *{{.Entity}}Repository) List(ctx context.Context, params *dbs.{{.ListParams}}) ([]*dbs.{{.ModelStruct}}, error) {
    return r.q.List{{.EntityPlural}}(ctx, *params)
}

func (r *{{.Entity}}Repository) GetByID(ctx context.Context, id string) (*dbs.{{.ModelStruct}}, error) {
    return r.q.Get{{.Entity}}ByID(ctx, id)
}

func (r *{{.Entity}}Repository) UpdateByID(ctx context.Context, params *dbs.{{.UpdateParams}}) (*dbs.{{.ModelStruct}}, error) {
    return r.q.Update{{.Entity}}ByID(ctx, *params)
}

func (r *{{.Entity}}Repository) DeleteByID(ctx context.Context, id string) error {
    return r.q.Delete{{.Entity}}ByID(ctx, id)
}

func (r *{{.Entity}}Repository) Count(ctx context.Context) (int64, error) {
    return r.q.Count{{.EntityPlural}}(ctx)
}
