package service

import (
    "context"

    "{{.ModulePath}}/internal/repository"
)

type I{{.Entity}}Service interface {
    Hello(ctx context.Context) error
}

type {{.Entity}}Service struct {
    repo repository.I{{.Entity}}Repository
}

func New{{.Entity}}Service(repo repository.I{{.Entity}}Repository) I{{.Entity}}Service {
    return &{{.Entity}}Service{repo: repo}
}

func (s *{{.Entity}}Service) Hello(ctx context.Context) error {
    // TODO: implement
    return s.repo.Hello(ctx)
}
