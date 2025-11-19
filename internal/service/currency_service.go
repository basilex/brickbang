package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"brickbang/internal/repository/dbs"
	"brickbang/internal/transfer"
)

type ICurrencyService interface {
	Create(ctx context.Context, req *transfer.CurrencyCreateRequest) (*dbs.Currency, error)
	GetByID(ctx context.Context, id string) (*dbs.Currency, error)
	GetByName(ctx context.Context, name string) (*dbs.Currency, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Currency, error)
	ListWithCountries(ctx context.Context, order string, limit, offset int32) ([]*dbs.ListCurrenciesWithCountriesRow, error)
	UpdateByID(ctx context.Context, id string, req *transfer.CurrencyUpdateRequest) (*dbs.Currency, error)
	DeleteByID(ctx context.Context, id string) (string, error)
}

type currencyService struct {
	queries *dbs.Queries
}

func NewCurrencyService(q *dbs.Queries) ICurrencyService {
	return &currencyService{queries: q}
}

// Create a new currency
func (rcv *currencyService) Create(ctx context.Context, req *transfer.CurrencyCreateRequest) (*dbs.Currency, error) {
	return rcv.queries.CreateCurrency(ctx, &dbs.CreateCurrencyParams{
		Name:    req.Name,
		Code:    req.Code,
		NumCode: req.NumCode,
		Symbol:  req.Symbol,
	})
}

// Get a currency by its ID
func (rcv *currencyService) GetByID(ctx context.Context, id string) (*dbs.Currency, error) {
	currency, err := rcv.queries.GetCurrencyByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("currency not found")
		}
		return nil, fmt.Errorf("failed to get currency by ID: %w", err)
	}

	return currency, nil
}

// Get a currency by its name
func (rcv *currencyService) GetByName(ctx context.Context, name string) (*dbs.Currency, error) {
	currency, err := rcv.queries.GetCurrencyByName(ctx, name)
	if err != nil {
		return nil, errors.New("currency not found")
	}

	return currency, nil
}

// Count total number of currencies
func (rcv *currencyService) Count(ctx context.Context) (int64, error) {
	return rcv.queries.CountCurrencies(ctx)
}

// List currencies with pagination and order
func (rcv *currencyService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Currency, error) {
	params := &dbs.ListCurrenciesParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	}

	return rcv.queries.ListCurrencies(ctx, params)
}

// List currencies with countries with pagination and order
func (rcv *currencyService) ListWithCountries(ctx context.Context, order string, limit, offset int32) ([]*dbs.ListCurrenciesWithCountriesRow, error) {
	params := &dbs.ListCurrenciesWithCountriesParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	}

	return rcv.queries.ListCurrenciesWithCountries(ctx, params)
}

// Update currency by ID
func (rcv *currencyService) UpdateByID(ctx context.Context, id string, req *transfer.CurrencyUpdateRequest) (*dbs.Currency, error) {
	params := &dbs.UpdateCurrencyByIDParams{
		ID:      id,
		Name:    req.Name,
		Code:    req.Code,
		NumCode: req.NumCode,
		Symbol:  req.Symbol,
	}

	return rcv.queries.UpdateCurrencyByID(ctx, params)
}

// Delete currency by ID
func (rcv *currencyService) DeleteByID(ctx context.Context, id string) (string, error) {
	id, err := rcv.queries.DeleteCurrencyByID(ctx, id)
	if err != nil {
		return "", err
	}

	return id, nil
}
