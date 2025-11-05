package repository

import (
    "context"
    "{{.ModulePath}}/storage/dbs"
)

type I{{.Entity}}Repository interface {
    Hello(ctx context.Context) error
}

type {{.Entity}}Repository struct {
    db *dbs.Queries
}

func New{{.Entity}}Repository(db *dbs.Queries) I{{.Entity}}Repository {
    return &{{.Entity}}Repository{db: db}
}

func (r *{{.Entity}}Repository) Hello(ctx context.Context) error {
    // TODO: implement
    return nil
}
