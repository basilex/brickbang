-- name: CurrencyNew :one
insert into currency (
	name, code, num_code, symbol
) values (
	@name, @code, @num_code, @symbol
) returning id, name, code, num_code, symbol, created_at, updated_at;

-- name: CurrencyCount :one
select count(*) from currency;

-- name: CurrencySelect :many
select *
  from currency c
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: CurrencyCountrySelect :many
select cr.id, cr.name, cr.code, cr.num_code, cr.symbol, cr.created_at, cr.updated_at,
  coalesce(
    json_agg(
      jsonb_build_object(
        'id',         cn.id,
        'name',       cn.name,
        'iso2',       cn.iso2,
        'iso3',       cn.iso3,
        'num_code',   cn.num_code, 
        'created_at', cn.created_at,
        'updated_at', cn.updated_at
      )
    ) filter (where cn.id is not null), '[]') as countries
  from currency cr
  left join country_currency cc on cr.id = cc.currency_id
  left join country cn ON cc.country_id = cn.id
 group by cr.id, cr.name
 order by @sql_order::text
 limit  @sql_limit offset @sql_offset;

-- name: CurrencySelectByID :one
select * from currency c where c.id = @id;

-- name: CurrencyUpdateByID :one
update currency
   set name = @name, code = @code, num_code = @num_code, symbol = @symbol
 where id = @id
       returning id, name, code, num_code, symbol, created_at, updated_at;

-- name: CurrencyDeleteByID :one
delete from currency c where c.id = @id returning id;
