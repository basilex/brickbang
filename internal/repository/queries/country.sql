-- name: CreateCountry :one
INSERT INTO country (
	name, iso2, iso3, num_code
) VALUES (
	@name, @iso2, @iso3, @num_code
)
RETURNING id, name, iso2, iso3, num_code, created_at, updated_at;

-- name: CountCountries :one
SELECT count(*) FROM country;

-- name: ListCountries :many
SELECT *
  FROM country c
 ORDER BY
  CASE WHEN @sql_order = 'asc' THEN name END ASC,
  CASE WHEN @sql_order = 'desc' THEN name END DESC,
    name ASC
 LIMIT @sql_limit OFFSET @sql_offset;

-- name: GetCountryByID :one
SELECT * FROM country c WHERE c.id = @id;

-- name: UpdateCountryByID :one
UPDATE country
   SET name = @name,
       iso2 = @iso2,
       iso3 = @iso3,
       num_code = @num_code
 WHERE id = @id
 RETURNING id, name, iso2, iso3, num_code, created_at, updated_at;

-- name: DeleteCountryByID :one
DELETE FROM country c WHERE c.id = @id RETURNING id;

-- name: ListCountriesWithCurrencies :many
SELECT cn.id,
       cn.name,
       cn.iso2,
       cn.iso3,
       cn.num_code,
       cn.created_at,
       cn.updated_at,
  COALESCE(
    JSON_AGG(
      JSONB_BUILD_OBJECT(
        'id',         cr.id,
        'name',       cr.name,
        'code',       cr.code,
        'num_code',   cr.num_code,
        'symbol',     cr.symbol,
        'created_at', cr.created_at,
        'updated_at', cr.updated_at
      )
    ) FILTER (WHERE cr.id IS NOT NULL), '[]'
  ) AS currencies
  FROM country cn
  LEFT JOIN country_currency cc ON cn.id = cc.country_id
  LEFT JOIN currency cr ON cc.currency_id = cr.id
 GROUP BY cn.id, cn.name
 ORDER BY
  CASE WHEN @sql_order = 'asc' THEN cn.name END ASC,
  CASE WHEN @sql_order = 'desc' THEN cn.name END DESC
 LIMIT  @sql_limit OFFSET @sql_offset;
