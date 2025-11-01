-- name: CountryNew :one
insert into country (
	name, iso2, iso3, num_code
) values (
	@name, @iso2, @iso3, @num_code
) returning id, name, iso2, iso3, num_code, created_at, updated_at;

-- name: CountryCount :one
select count(*) from country;

-- name: CountrySelect :many
select *
  from country c
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: CountryCurrencySelect :many
select cn.id, cn.name, cn.iso2, cn.iso3, cn.num_code, cn.created_at, cn.updated_at,
  coalesce(
    json_agg(
      jsonb_build_object(
        'id',         cr.id,
        'name',       cr.name,
        'code',       cr.code,
        'num_code',   cr.num_code,
        'symbol',     cr.symbol,
        'created_at', cr.created_at,
        'updated_at', cr.updated_at
      )
    ) filter (where cr.id is not null), '[]') as currencies
  from country cn
  left join country_currency cc on cn.id = cc.country_id
  left join currency cr ON cc.currency_id = cr.id
 group by cn.id, cn.name
 order by @sql_order::text
 limit  @sql_limit offset @sql_offset;

-- name: CountrySelectByID :one
select * from country c where c.id = @id;

-- name: CountryUpdateByID :one
update country
   set name = @name, iso2 = @iso2, iso3 = @iso3, num_code = @num_code
 where id = @id
       returning id, name, iso2, iso3, num_code, created_at, updated_at;

-- name: CountryDeleteByID :one
delete from country c where c.id = @id returning id;
