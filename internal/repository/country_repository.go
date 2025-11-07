package repository

import (
    "context"

    "brickbang/storage/dbs"
)

type ICountryRepository interface {
    Create(ctx context.Context, params *dbs.CreateCountryParams) (*dbs., error)
    List(ctx context.Context, params *dbs.ListCountriesWithCurrenciesParams) ([]*dbs., error)
    GetByID(ctx context.Context, id string) (*dbs., error)
    UpdateByID(ctx context.Context, params *dbs.UpdateCountryByIDParams) (*dbs., error)
    DeleteByID(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}

type CountryRepository struct {
    q *dbs.Queries
}

func NewCountryRepository(q *dbs.Queries) ICountryRepository {
    return &CountryRepository{q: q}
}

func (r *CountryRepository) Create(ctx context.Context, params *dbs.CreateCountryParams) (*dbs., error) {
    return r.q.CreateCountry(ctx, *params)
}

func (r *CountryRepository) List(ctx context.Context, params *dbs.ListCountriesWithCurrenciesParams) ([]*dbs., error) {
    return r.q.ListCountries(ctx, *params)
}

func (r *CountryRepository) GetByID(ctx context.Context, id string) (*dbs., error) {
    return r.q.GetCountryByID(ctx, id)
}

func (r *CountryRepository) UpdateByID(ctx context.Context, params *dbs.UpdateCountryByIDParams) (*dbs., error) {
    return r.q.UpdateCountryByID(ctx, *params)
}

func (r *CountryRepository) DeleteByID(ctx context.Context, id string) error {
    return r.q.DeleteCountryByID(ctx, id)
}

func (r *CountryRepository) Count(ctx context.Context) (int64, error) {
    return r.q.CountCountries(ctx)
}
