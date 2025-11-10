drop procedure if exists delete_if_exists(varchar);

drop function if exists trigger_updated_at();

drop function if exists public.xid_counter(public.xid);
drop function if exists public.xid_pid(public.xid);
drop function if exists public.xid_machine(public.xid);
drop function if exists public.xid_time(public.xid);
drop function if exists xid(timestamptz);
drop function if exists public.xid_decode(public.xid);
drop function if exists public.xid_encode(int[]);
drop function if exists public._xid_machine_id();

drop sequence if exists public.xid_serial;
drop domain if exists public.xid;

drop extension if exists "pgcrypto";
drop extension if exists "hstore";
