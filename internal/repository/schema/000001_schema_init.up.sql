--
-- BSPDB postgresql schema
--
create extension if not exists "hstore";
create extension if not exists "pgcrypto";
--
-- Generate XID as Postgres unique value
--
do $$ begin
    create domain public.xid as char(20) check (value ~ '^[a-v0-9]{20}$');
exception
    when duplicate_object then null;
end $$;

create sequence if not exists public.xid_serial minvalue 0 maxvalue 16777215 cycle; --  ((255<<16) + (255<<8) + 255))

select setval('xid_serial', (random() * 16777215)::int); --  ((255<<16) + (255<<8) + 255))

create or replace function public._xid_machine_id()
    returns int
    language plpgsql
    immutable
as
$$
declare
begin
    return (select system_identifier & 16777215 from pg_control_system());
end
$$;

create or replace function public.xid_encode(_id int[])
    returns public.xid
    language plpgsql
as
$$
declare
    _encoding char(1)[] = '{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v}';
begin
    return _encoding[1 + (_id[1] >> 3)]
               || _encoding[1 + ((_id[2] >> 6) & 31 | (_id[1] << 2) & 31)]
               || _encoding[1 + ((_id[2] >> 1) & 31)]
               || _encoding[1 + ((_id[3] >> 4) & 31 | (_id[2] << 4) & 31)]
               || _encoding[1 + (_id[4] >> 7 | (_id[3] << 1) & 31)]
               || _encoding[1 + ((_id[4] >> 2) & 31)]
               || _encoding[1 + (_id[5] >> 5 | (_id[4] << 3) & 31)]
               || _encoding[1 + (_id[5] & 31)]
               || _encoding[1 + (_id[6] >> 3)]
               || _encoding[1 + ((_id[7] >> 6) & 31 | (_id[6] << 2) & 31)]
               || _encoding[1 + ((_id[7] >> 1) & 31)]
               || _encoding[1 + ((_id[8] >> 4) & 31 | (_id[7] << 4) & 31)]
               || _encoding[1 + (_id[9] >> 7 | (_id[8] << 1) & 31)]
               || _encoding[1 + ((_id[9] >> 2) & 31)]
               || _encoding[1 + ((_id[10] >> 5) | (_id[9] << 3) & 31)]
               || _encoding[1 + (_id[10] & 31)]
               || _encoding[1 + (_id[11] >> 3)]
               || _encoding[1 + ((_id[12] >> 6) & 31 | (_id[11] << 2) & 31)]
               || _encoding[1 + ((_id[12] >> 1) & 31)]
        || _encoding[1 + ((_id[12] << 4) & 31)];
end;
$$;

create or replace function public.xid_decode(_xid public.xid)
    returns int[]
    language plpgsql
as
$$
declare
    _dec int[] = '{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}';
    _b   bytea;
begin
    _b := _xid::bytea;
    return array [
            ((_dec[get_byte(_b, 0)] << 3) | (_dec[get_byte(_b, 1)] >> 2)) & 255,
            ((_dec[get_byte(_b, 1)] << 6) | (_dec[get_byte(_b, 2)] << 1) | (_dec[get_byte(_b, 3)] >> 4)) & 255,
            ((_dec[get_byte(_b, 3)] << 4) | (_dec[get_byte(_b, 4)] >> 1)) & 255,
            ((_dec[get_byte(_b, 4)] << 7) | (_dec[get_byte(_b, 5)] << 2) | (_dec[get_byte(_b, 6)] >> 3)) & 255,
        ((_dec[get_byte(_b, 6)] << 5) | (_dec[get_byte(_b, 7)])) & 255,
            ((_dec[get_byte(_b, 8)] << 3) | (_dec[get_byte(_b, 9)] >> 2)) & 255,
            ((_dec[get_byte(_b, 9)] << 6) | (_dec[get_byte(_b, 10)] << 1) | (_dec[get_byte(_b, 11)] >> 4)) & 255,
            ((_dec[get_byte(_b, 11)] << 4) | (_dec[get_byte(_b, 12)] >> 1)) & 255,
            ((_dec[get_byte(_b, 12)] << 7) | (_dec[get_byte(_b, 13)] << 2) | (_dec[get_byte(_b, 14)] >> 3)) & 255,
        ((_dec[get_byte(_b, 14)] << 5) | (_dec[get_byte(_b, 15)])) & 255,
            ((_dec[get_byte(_b, 16)] << 3) | (_dec[get_byte(_b, 17)] >> 2)) & 255,
            ((_dec[get_byte(_b, 17)] << 6) | (_dec[get_byte(_b, 18)] << 1) | (_dec[get_byte(_b, 19)] >> 4)) & 255
        ];
end;
$$;

create or replace function xid(_at timestamptz default current_timestamp)
    returns public.xid
    language plpgsql
as
$$
declare
    _t int;
    _m int;
    _p int;
    _c int;
begin
    _t := floor(extract(epoch from _at));
    _m := _xid_machine_id();
    _p := pg_backend_pid();
    _c := nextval('xid_serial')::int;

    return public.xid_encode(array [
        (_t >> 24) & 255, (_t >> 16) & 255, (_t >> 8) & 255 , _t & 255,
        (_m >> 16) & 255, (_m >> 8) & 255 , _m & 255,
        (_p >> 8) & 255, _p & 255,
        (_c >> 16) & 255, (_c >> 8) & 255 , _c & 255
        ]);
end;
$$;

create or replace function public.xid_time(_xid public.xid)
    returns timestamptz
    language plpgsql
as
$$
declare
    _id int[];
begin
    _id := public.xid_decode(_xid);
    return to_timestamp((_id[1] << 24)::bigint + (_id[2] << 16) + (_id[3] << 8) + (_id[4]));
end;
$$;

create or replace function public.xid_machine(_xid public.xid)
    returns int[]
    language plpgsql
as
$$
declare
    _id int[];
begin
    _id := public.xid_decode(_xid);
    return array [_id[5], _id[6], _id[7]];
end;
$$;

create or replace function public.xid_pid(_xid public.xid)
    returns int
    language plpgsql
as
$$
declare
    _id int[];
begin
    _id := public.xid_decode(_xid);
    return (_id[8] << 8) + (_id[9]);
end;
$$;

create or replace function public.xid_counter(_xid public.xid)
    returns int
    language plpgsql
as
$$
declare
    _id int[];
begin
    _id := public.xid_decode(_xid);
    return (_id[10] << 16) + (_id[11] << 8) + (_id[12]);
end;
$$;

--
-- Trigger function for updating updated_at field
--
create or replace function trigger_updated_at()
returns trigger as $$
begin
   if row(new.*) is distinct from row(old.*) then
      new.updated_at = timezone('utc', now());
      return new;
   else
      return old;
   end if;
end;
$$ language 'plpgsql';
--
-- Delete from table if it exists
--
create or replace procedure delete_if_exists(in tbl varchar) as $$
declare
    is_exists boolean;
begin
   select exists (
      select from pg_tables
	    where schemaname = 'public'
	      and tablename = tbl
   ) into is_exists;
   if is_exists then
      raise notice '>>> deleting data from [%] table...', tbl;
	   execute 'delete from ' || tbl;
   end if;
end $$ language plpgsql;
--
-- eof
--
