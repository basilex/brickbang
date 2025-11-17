-- name: CreateCurrency :one
INSERT INTO currency (
	name, code, num_code, symbol
) VALUES (
	@name, @code, @num_code, @symbol
)
RETURNING id, name, code, num_code, symbol, created_at, updated_at;

-- name: CountCurrencies :one
SELECT count(*) FROM currency;

-- name: ListCurrencies :many
SELECT *
  FROM currency c
 ORDER BY
    CASE WHEN @sql_order = 'asc' THEN name END ASC,
    CASE WHEN @sql_order = 'desc' THEN name END DESC
 LIMIT @sql_limit OFFSET @sql_offset;

-- name: GetCurrencyByID :one
SELECT * FROM currency c WHERE c.id = @id;

-- name: UpdateCurrencyByID :one
UPDATE currency
   SET name = @name, code = @code, num_code = @num_code, symbol = @symbol
 WHERE id = @id
       RETURNING id, name, code, num_code, symbol, created_at, updated_at;

-- name: DeleteCurrencyByID :one
DELETE FROM currency c WHERE c.id = @id RETURNING id;

-- name: ListCurrenciesWithCountries :many
SELECT cr.id,
       cr.name,
       cr.code,
       cr.num_code,
       cr.symbol,
       cr.created_at,
       cr.updated_at,
  COALESCE(
    JSON_AGG(
      JSONB_BUILD_OBJECT(
        'id',         cn.id,
        'name',       cn.name,
        'iso2',       cn.iso2,
        'iso3',       cn.iso3,
        'num_code',   cn.num_code,
        'created_at', cn.created_at,
        'updated_at', cn.updated_at
      )
    ) FILTER (WHERE cn.id IS NOT NULL), '[]'
  ) AS countries
  FROM currency cr
  LEFT JOIN country_currency cc ON cr.id = cc.currency_id
  LEFT JOIN country cn ON cc.country_id = cn.id
 GROUP BY cr.id, cr.name
 ORDER BY
    CASE WHEN @sql_order = 'asc' THEN cr.name END ASC,
    CASE WHEN @sql_order = 'desc' THEN cr.name END DESC
 LIMIT @sql_limit OFFSET @sql_offset;
