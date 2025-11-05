package service

import (
	"context"

	"brickbang/internal/repository"
	"brickbang/storage/dbs"
)

type ICountryService interface {
	GetAll(ctx context.Context, order string, offset, limit int32) ([]*dbs.Country, error)
	GetAllWithCurrencies(ctx context.Context, order string, offset, limit int32) ([]*dbs.CountryCurrencySelectRow, error)
	GetByID(ctx context.Context, id string) (*dbs.Country, error)
	Create(ctx context.Context, dto *dbs.CountryNewParams) (*dbs.Country, error)
	Update(ctx context.Context, dto *dbs.CountryUpdateByIDParams) (*dbs.Country, error)
	Delete(ctx context.Context, id string) (string, error)
	Count(ctx context.Context) (int64, error)
}

type countryService struct {
	repo repository.ICountryRepository
}

func NewCountryService(repo repository.ICountryRepository) ICountryService {
	return &countryService{repo: repo}
}

func (s *countryService) GetAll(ctx context.Context, order string, offset, limit int32) ([]*dbs.Country, error) {
	params := &dbs.CountrySelectParams{
		SqlOrder:  order,
		SqlOffset: offset,
		SqlLimit:  limit,
	}
	return s.repo.Select(ctx, params)
}

func (s *countryService) GetAllWithCurrencies(ctx context.Context, order string, offset, limit int32) ([]*dbs.CountryCurrencySelectRow, error) {
	params := &dbs.CountryCurrencySelectParams{
		SqlOrder:  order,
		SqlOffset: offset,
		SqlLimit:  limit,
	}
	return s.repo.SelectWithCurrencies(ctx, params)
}

func (s *countryService) GetByID(ctx context.Context, id string) (*dbs.Country, error) {
	return s.repo.SelectByID(ctx, id)
}

func (s *countryService) Create(ctx context.Context, dto *dbs.CountryNewParams) (*dbs.Country, error) {
	return s.repo.Create(ctx, dto)
}

func (s *countryService) Update(ctx context.Context, dto *dbs.CountryUpdateByIDParams) (*dbs.Country, error) {
	return s.repo.Update(ctx, dto)
}

func (s *countryService) Delete(ctx context.Context, id string) (string, error) {
	return s.repo.Delete(ctx, id)
}

func (s *countryService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
