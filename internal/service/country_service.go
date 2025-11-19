package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"brickbang/internal/repository/dbs"
	"brickbang/internal/transfer"
)

type ICountryService interface {
	Create(ctx context.Context, req *transfer.CountryCreateRequest) (*dbs.Country, error)
	GetByID(ctx context.Context, id string) (*dbs.Country, error)
	GetByName(ctx context.Context, name string) (*dbs.Country, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Country, error)
	ListWithCurrencies(ctx context.Context, order string, limit, offset int32) ([]*dbs.ListCountriesWithCurrenciesRow, error)
	UpdateByID(ctx context.Context, id string, req *transfer.CountryUpdateRequest) (*dbs.Country, error)
	DeleteByID(ctx context.Context, id string) (string, error)
}

type countryService struct {
	queries *dbs.Queries
}

func NewCountryService(q *dbs.Queries) ICountryService {
	return &countryService{queries: q}
}

// Create a new country
func (rcv *countryService) Create(ctx context.Context, req *transfer.CountryCreateRequest) (*dbs.Country, error) {
	return rcv.queries.CreateCountry(ctx, &dbs.CreateCountryParams{
		Name:    req.Name,
		Iso2:    req.Iso2,
		Iso3:    req.Iso3,
		NumCode: req.NumCode,
	})
}

// Get a country by its ID
func (rcv *countryService) GetByID(ctx context.Context, id string) (*dbs.Country, error) {
	country, err := rcv.queries.GetCountryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("country not found")
		}
		return nil, fmt.Errorf("failed to get country by ID: %w", err)
	}

	return country, nil
}

// Get a country by its name
func (rcv *countryService) GetByName(ctx context.Context, name string) (*dbs.Country, error) {
	country, err := rcv.queries.GetCountryByName(ctx, name)
	if err != nil {
		return nil, errors.New("country not found")
	}

	return country, nil
}

// Count total number of countries
func (rcv *countryService) Count(ctx context.Context) (int64, error) {
	return rcv.queries.CountCountries(ctx)
}

// List countries with pagination and order
func (rcv *countryService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Country, error) {
	params := &dbs.ListCountriesParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	}

	return rcv.queries.ListCountries(ctx, params)
}

// List countries with currencies with pagination and order
func (rcv *countryService) ListWithCurrencies(ctx context.Context, order string, limit, offset int32) ([]*dbs.ListCountriesWithCurrenciesRow, error) {
	params := &dbs.ListCountriesWithCurrenciesParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	}

	return rcv.queries.ListCountriesWithCurrencies(ctx, params)
}

// Update country by ID
func (rcv *countryService) UpdateByID(ctx context.Context, id string, req *transfer.CountryUpdateRequest) (*dbs.Country, error) {
	params := &dbs.UpdateCountryByIDParams{
		ID:      id,
		Name:    req.Name,
		Iso2:    req.Iso2,
		Iso3:    req.Iso3,
		NumCode: req.NumCode,
	}

	return rcv.queries.UpdateCountryByID(ctx, params)
}

// Delete country by ID
func (rcv *countryService) DeleteByID(ctx context.Context, id string) (string, error) {
	id, err := rcv.queries.DeleteCountryByID(ctx, id)
	if err != nil {
		return "", err
	}

	return id, nil
}
