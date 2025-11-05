package repository

import (
	"context"

	"brickbang/storage/dbs"
)

type ICountryRepository interface {
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context, params *dbs.CountrySelectParams) ([]*dbs.Country, error)
	SelectByID(ctx context.Context, id string) (*dbs.Country, error)
	SelectWithCurrencies(ctx context.Context, params *dbs.CountryCurrencySelectParams) ([]*dbs.CountryCurrencySelectRow, error)
	Create(ctx context.Context, params *dbs.CountryNewParams) (*dbs.Country, error)
	Update(ctx context.Context, params *dbs.CountryUpdateByIDParams) (*dbs.Country, error)
	Delete(ctx context.Context, id string) (string, error)
}

type countryRepository struct {
	q *dbs.Queries
}

func NewCountryRepository(q *dbs.Queries) ICountryRepository {
	return &countryRepository{q: q}
}

func (r *countryRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountryCount(ctx)
}

func (r *countryRepository) Select(ctx context.Context, params *dbs.CountrySelectParams) ([]*dbs.Country, error) {
	return r.q.CountrySelect(ctx, params)
}

func (r *countryRepository) SelectByID(ctx context.Context, id string) (*dbs.Country, error) {
	return r.q.CountrySelectByID(ctx, id)
}

func (r *countryRepository) SelectWithCurrencies(ctx context.Context, params *dbs.CountryCurrencySelectParams) ([]*dbs.CountryCurrencySelectRow, error) {
	return r.q.CountryCurrencySelect(ctx, params)
}

func (r *countryRepository) Create(ctx context.Context, params *dbs.CountryNewParams) (*dbs.Country, error) {
	return r.q.CountryNew(ctx, params)
}

func (r *countryRepository) Update(ctx context.Context, params *dbs.CountryUpdateByIDParams) (*dbs.Country, error) {
	return r.q.CountryUpdateByID(ctx, params)
}

func (r *countryRepository) Delete(ctx context.Context, id string) (string, error) {
	return r.q.CountryDeleteByID(ctx, id)
}
