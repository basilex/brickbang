--
-- Entity user
--
create table users (
  id              varchar(32)     not null default xid() primary key,
  username        varchar(64)     not null unique,
  password        varchar(255)    not null,

  is_blocked      bool            not null default false,
  blocked_at      timestamp       not null default '1000-01-01'::timestamp,

  is_checked      bool            not null default false,
  checked_at      timestamp       not null default '1000-01-01'::timestamp,

  visited_at      timestamp       not null default '1000-01-01'::timestamp,
  created_at      timestamp       not null default timezone('utc', now()),
  updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create trigger users_updated_at
	before update on users for each row
	execute procedure trigger_updated_at();

insert into users(username, password, is_checked, checked_at) values
  ('sys', crypt('passw!rd', gen_salt('bf', 12)), true, timezone('utc', now())),
  ('admin', crypt('passw=rd', gen_salt('bf', 12)), true, timezone('utc', now())),
  ('basilex', crypt('passw0rd', gen_salt('bf', 12)), true, timezone('utc', now()));
--
-- Entity roles
--
create table roles (
  id              varchar(32)     not null default xid() primary key,
  pid             varchar(32)     not null default xid() unique,
  name            varchar(255)    not null unique,
  description     text            not null default '',
  created_at      timestamp       not null default timezone('utc', now()),
  updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create trigger users_updated_at
	before update on roles for each row
	execute procedure trigger_updated_at();

insert into roles(name, description) values
  ('role_sys', 'Role with superuser permissions'),
  ('role_admin', 'Role with administrative permissions'),
  ('role_editor', 'Role with content editing permissions'),
  ('role_manager', 'Role with management permissions'),
  ('role_support', 'Role with support permissions'),
  ('role_reporter', 'Role with reporting permissions'),
  ('role_financier', 'Role with financial permissions'),
  ('role_customer', 'Role with customer permissions');
--
-- Entity user_roles
--
create table user_roles (
  id              varchar(32)     not null default xid() primary key,
  user_id         varchar(32)     not null references users(id) on delete cascade,
  role_id         varchar(32)     not null references roles(id) on delete cascade,
  created_at      timestamp       not null default timezone('utc', now())
);

create unique index user_roles_pair_unq on user_roles(user_id, role_id);

create or replace function notify_user_roles() returns trigger as $$
declare
  payload json;
  uid text;
  rid text;
begin
  if (TG_OP = 'INSERT') then
    uid := new.user_id;
    rid := new.role_id;
  elsif (TG_OP = 'DELETE') then
    uid := old.user_id;
    rid := old.role_id;
  else
    uid := new.user_id;
    rid := new.role_id;
  end if;

  payload := json_build_object(
    'object', 'user_roles',
    'op', TG_OP,
    'user_id', uid,
    'role_id', rid
  );
  perform pg_notify('rbac_updates', payload::text);
  return case when TG_OP = 'DELETE' then old else new end;
end;
$$ language plpgsql;

create trigger user_roles_change
  after insert or update or delete on user_roles
  for each row execute function notify_user_roles();

do $$
declare
    m record;
begin
  for m in
    select * from (values
      ('sys',     'role_sys'),
      ('admin',   'role_admin'),
      ('basilex', 'role_admin')
    ) as t(username, role_name)
  loop
    insert into user_roles(user_id, role_id)
      select u.id, r.id
        from users u, roles r
        where u.username = m.username
          and r.name = m.role_name;
  end loop;
end $$;
--
-- Entity grants
--
create table grants (
  id              varchar(32)     not null default xid() primary key,
  pid             varchar(32)     not null default xid() unique,
  code            varchar(255)    not null unique,
  description     text            not null default '',
  created_at      timestamp       not null default timezone('utc', now()),
  updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create trigger grants_updated_at
  before update on grants for each row
  execute procedure trigger_updated_at();

create or replace function notify_grants() returns trigger as $$
declare
  payload json;
  gid text;
  code text;
begin
  if (TG_OP = 'INSERT') then
    gid := new.id;
    code := new.code;
  elsif (TG_OP = 'DELETE') then
    gid := old.id;
    code := old.code;
  else
    gid := new.id;
    code := new.code;
  end if;

  payload := json_build_object(
    'object', 'grants',
    'op', TG_OP,
    'grant_id', gid,
    'code', code
  );
  perform pg_notify('rbac_updates', payload::text);
  return case when TG_OP = 'DELETE' then old else new end;
end;
$$ language plpgsql;

create trigger grants_change
  after insert or update or delete on grants
  for each row execute function notify_grants();

do $$
begin
  insert into grants(code, description) values
    ('sys:*:*', 'Superuser all permissions'),

    ('admin:users:*', 'Manage users'),
    ('admin:roles:*', 'Manage roles'),
    ('admin:grants:*', 'Manage grants'),
    ('admin:contacts:*', 'Manage contacts'),
    ('admin:profiles:*', 'Manage profiles'),
    ('admin:sessions:*', 'Manage sessions'),
    ('admin:countries:*', 'Manage country data'),
    ('admin:currencies:*', 'Manage currency data'),
    ('admin:settings:*', 'Manage system settings'),

    ('tenant:users:read', 'Read user info'),
    ('tenant:users:update', 'Update user info'),
    ('tenant:roles:assign', 'Assign roles to users'),
    ('tenant:settings:read', 'Read system settings'),
    ('tenant:settings:update', 'Update system settings'),

    ('tenant:product:create', 'Create products'),
    ('tenant:product:read', 'Read products'),
    ('tenant:product:update', 'Update product info'),
    ('tenant:product:delete', 'Delete product at all'),

    ('tenant:order:create', 'Create orders'),
    ('tenant:order:read', 'Read orders'),
    ('tenant:order:update', 'Update order info'),
    ('tenant:order:delete', 'Delete orders at all');
end $$;
--
-- Entity role_grants
--
create table role_grants (
  id              varchar(32)     not null default xid() primary key,
  role_id         varchar(32)     not null references roles(id) on delete cascade,
  grant_id        varchar(32)     not null references grants(id) on delete cascade,
  created_at      timestamp       not null default timezone('utc', now())
);

create unique index role_grants_pair_unq on role_grants(role_id, grant_id);

create or replace function notify_role_grants() returns trigger as $$
declare
  payload json;
  rid text;
  gid text;
begin
  if (TG_OP = 'INSERT') then
    rid := new.role_id;
    gid := new.grant_id;
  elsif (TG_OP = 'DELETE') then
    rid := old.role_id;
    gid := old.grant_id;
  else
    -- UPDATE: handle like INSERT (values on NEW)
    rid := new.role_id;
    gid := new.grant_id;
  end if;

  payload := json_build_object(
    'object', 'role_grants',
    'op', TG_OP,
    'role_id', rid,
    'grant_id', gid
  );
  perform pg_notify('rbac_updates', payload::text);
  return case when TG_OP = 'DELETE' then old else new end;
end;
$$ language plpgsql;

create trigger role_grants_change
  after insert or update or delete on role_grants
  for each row execute function notify_role_grants();

do $$
declare
  v_role_id role_grants.role_id%type;
begin
  -- SYS ROLE — получает *все* sys:*:* права
  select id into v_role_id from roles where name = 'role_sys';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code like 'sys:%'
  on conflict do nothing;

  -- ADMIN ROLE — получает все admin:*:* права
  select id into v_role_id from roles where name = 'role_admin';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code like 'admin:%'
  on conflict do nothing;

  -- EDITOR ROLE (пример: доступ только к чтению/обновлению каталогов)
  select id into v_role_id from roles where name = 'role_editor';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code in (
    'tenant:product:read',
    'tenant:product:update',
    'tenant:order:read',
    'tenant:order:update'
  )
  on conflict do nothing;

  -- MANAGER ROLE (пример: товары/заказы create+read+update)
  select id into v_role_id from roles where name = 'role_manager';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code in (
    'tenant:product:create',
    'tenant:product:read',
    'tenant:product:update',
    'tenant:order:create',
    'tenant:order:read',
    'tenant:order:update'
  )
  on conflict do nothing;

  -- SUPPORT ROLE — просмотр данных пользователей и заказов
  select id into v_role_id from roles where name = 'role_support';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code in (
    'tenant:users:read',
    'tenant:product:read',
    'tenant:order:read'
  )
  on conflict do nothing;

  -- REPORTER ROLE — только read, любые данные
  select id into v_role_id from roles where name = 'role_reporter';
  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code like '%:read'
  on conflict do nothing;

  -- FINANCIER ROLE — доступ к заказам + настройкам
  select id into v_role_id from roles where name = 'role_financier';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code in (
    'tenant:order:read',
    'tenant:order:update',
    'tenant:settings:read'
  )
  on conflict do nothing;

  -- CUSTOMER ROLE — минимальные права
  select id into v_role_id from roles where name = 'role_customer';

  insert into role_grants(role_id, grant_id)
  select v_role_id, g.id
    from grants g
   where g.code in (
    'tenant:users:read',
    'tenant:product:read'
  )
  on conflict do nothing;
end $$;
--
-- Entity session
--
create table session (
  id              varchar(32)     not null default xid() primary key,
  user_id         varchar(32)     not null references users(id) on delete cascade,

  access_token    varchar(1024)   not null,
  refresh_token   varchar(1024)   not null,

  access_exp      timestamp       not null,
  refresh_exp     timestamp       not null,

  access_jti      varchar(256)   not null default '',
  refresh_jti     varchar(256)   not null default '',

  access_status   varchar(32)     not null default 'valid' check(access_status in ('valid','expired','revoked')),
  refresh_status  varchar(32)     not null default 'valid' check(refresh_status in ('valid','expired','revoked')),

  ip_address      varchar(64)     not null default '-',
  user_agent      varchar(512)    not null default '-',

  created_at      timestamp       not null default timezone('utc', now()),
  updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create trigger session_updated_at
  before update on session for each row
  execute procedure trigger_updated_at();
--
-- Entity profile
--
create table profile (
  id              varchar(32)     not null default xid() primary key,
  user_id         varchar(32)     not null unique references users(id) on delete cascade,
  firstname       varchar(255)    not null,
  lastname        varchar(255)    not null,
  gender          varchar(32)     not null default 'unknown' check(gender in('unknown', 'male', 'female')),
  birthday        date            not null default '1000-01-01'::date,
  avatar_url      varchar(255)    not null default '/assets/images/person.jpg',
  enable_2fa      bool            not null default false,
  secret_2fa      varchar(255),
  created_at      timestamp       not null default timezone('utc', now()),
  updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create trigger profile_updated_at
	before update on profile for each row
	execute procedure trigger_updated_at();

do $$
declare
  v_user_id profile.id%type;
begin
  select id into v_user_id from users where username = 'sys';
  insert into profile(user_id, firstname, lastname)
    values(v_user_id, 'Superuser', 'Administrator');

  select id into v_user_id from users where username = 'admin';
  insert into profile(user_id, firstname, lastname)
    values(v_user_id, 'System', 'Administrator');

  select id into v_user_id from users where username = 'basilex';
  insert into profile(user_id, firstname, lastname, gender, birthday)
    values(v_user_id, 'Alexander', 'Vasilenko', 'male', '1965-04-03'::date);
end $$;
--
-- Entity contact
--
create table contact (
  id              varchar(32)     not null default xid() primary key,
  user_id         varchar(32)     not null references users(id) on delete cascade,
  class           varchar(32)     not null check(class in ('email', 'phone', 'mobile', 'telegram', 'viber', 'signal', 'other')),
  content         varchar(255)    not null,
  created_at      timestamp       not null default timezone('utc', now()),
  updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create index contact_user_id on contact(user_id);

create unique index contact_class_content_unq on contact(user_id, class, content);

create unique index contact_class_email_unq on contact(content) where class = 'email';
create unique index contact_class_phone_unq on contact(content) where class = 'phone';
create unique index contact_class_mobile_unq on contact(content) where class = 'mobile';
create unique index contact_class_telegram_unq on contact(content) where class = 'telegram';
create unique index contact_class_viber_unq on contact(content) where class = 'viber';
create unique index contact_class_signal_unq on contact(content) where class = 'signal';

create trigger contact_updated_at
	before update on contact for each row
	execute procedure trigger_updated_at();

do $$
declare
  v_user_id users.id%type;
begin
  --
  -- sys contacts
  --
  select id into v_user_id from users where username = 'sys';
  insert into contact(user_id, class, content) values(v_user_id, 'email', 'sys@brickbang.com');
  --
  -- admin contacts
  --
  select id into v_user_id from users where username = 'admin';
  insert into contact(user_id, class, content) values(v_user_id, 'email', 'admin@brickbang.com');
  --
  -- basilex contacts
  --
  select id into v_user_id from users where username = 'basilex';

  insert into contact(user_id, class, content) values(v_user_id, 'email', 'alexander.vasilenko@gmail.com');
  insert into contact(user_id, class, content) values(v_user_id, 'email', 'alexander.vasilenko@icloud.com');
  insert into contact(user_id, class, content) values(v_user_id, 'mobile', '+380952066922');
  insert into contact(user_id, class, content) values(v_user_id, 'telegram', '@Basilex');
end $$;
