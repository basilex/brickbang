package service

import (
	"context"

	"brickbang/internal/mapper"
	"brickbang/internal/repository"
)

type I{{.Entity}}Service interface {
	Create(ctx context.Context, req *mapper.{{.Entity}}CreateRequest) (*mapper.{{.Entity}}Response, error)
	Update(ctx context.Context, req *mapper.{{.Entity}}UpdateRequest) (*mapper.{{.Entity}}Response, error)
	Get(ctx context.Context, id string) (*mapper.{{.Entity}}Response, error)
	List(ctx context.Context) ([]*mapper.{{.Entity}}Response, error)
	Delete(ctx context.Context, id string) error
}

type {{.Entity}}Service struct {
	repo repository.I{{.Entity}}Repository
}

func New{{.Entity}}Service(repo repository.I{{.Entity}}Repository) I{{.Entity}}Service {
	return &{{.Entity}}Service{repo: repo}
}

func (s *{{.Entity}}Service) Create(ctx context.Context, req *mapper.{{.Entity}}CreateRequest) (*mapper.{{.Entity}}Response, error) {
	params := mapper.To{{.Entity}}NewParams(req)
	row, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.To{{.Entity}}Response(row), nil
}

func (s *{{.Entity}}Service) Update(ctx context.Context, req *mapper.{{.Entity}}UpdateRequest) (*mapper.{{.Entity}}Response, error) {
	params := mapper.To{{.Entity}}UpdateParams(req.ID, req)
	row, err := s.repo.Update(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.To{{.Entity}}Response(row), nil
}

func (s *{{.Entity}}Service) Get(ctx context.Context, id string) (*mapper.{{.Entity}}Response, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.To{{.Entity}}Response(row), nil
}

func (s *{{.Entity}}Service) List(ctx context.Context) ([]*mapper.{{.Entity}}Response, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.To{{.Entity}}ResponseList(rows), nil
}

func (s *{{.Entity}}Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
