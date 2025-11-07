package service

import (
    "context"

    "brickbang/storage/dbs"
)

// Интерфейс по конвенции I...
type ICountryService interface {
    Create(ctx context.Context, params dbs.CreateCountryParams) (*dbs.Country, error)
    List(ctx context.Context, params dbs.ListCountriesWithCurrenciesParams) ([]*dbs.Country, error)
    GetByID(ctx context.Context, id string) (*dbs.Country, error)
    
    UpdateByID(ctx context.Context, params dbs.UpdateCountryByIDParams) (*dbs.Country, error)
    DeleteByID(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}

// Реализация сервиса
type CountryService struct {
    q *dbs.Queries
}

func NewCountryService(q *dbs.Queries) ICountryService {
    return &CountryService{q: q}
}


func (s *CountryService) Create(ctx context.Context, params dbs.CreateCountryParams) (*dbs.Country, error) {
    return s.q.CreateCountry(ctx, params)
}



func (s *CountryService) List(ctx context.Context, params dbs.ListCountriesWithCurrenciesParams) ([]*dbs.Country, error) {
    return s.q.ListCountrys(ctx, params)
}



func (s *CountryService) GetByID(ctx context.Context, id string) (*dbs.Country, error) {
    return s.q.GetCountryByID(ctx, id)
}





func (s *CountryService) UpdateByID(ctx context.Context, params dbs.UpdateCountryByIDParams) (*dbs.Country, error) {
    return s.q.UpdateCountryByID(ctx, params)
}



func (s *CountryService) DeleteByID(ctx context.Context, id string) error {
    return s.q.DeleteCountryByID(ctx, id)
}



func (s *CountryService) Count(ctx context.Context) (int64, error) {
    return s.q.CountCountry(ctx)
}

