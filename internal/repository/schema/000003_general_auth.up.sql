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
    name            varchar(255)    not null unique,
    created_at      timestamp       not null default timezone('utc', now()),
    updated_at      timestamp       not null default '1000-01-01'::timestamp
);

create trigger users_updated_at
	before update on roles for each row
	execute procedure trigger_updated_at();

insert into roles(name) values
    ('SYS'),
    ('ADMIN'),
    ('EDITOR'),
    ('AUDITOR'),
    ('MANAGER'),
    ('REPORTER'),
    ('FINANCIER');
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

do $$
declare
    v_user_id user_roles.user_id%type;
    v_role_id user_roles.role_id%type;
begin
    select id into v_user_id from users where username = 'sys';
    select id into v_role_id from roles where name = 'SYS';
    insert into user_roles(user_id, role_id) values(v_user_id, v_role_id);

    select id into v_user_id from users where username = 'admin';
    select id into v_role_id from roles where name = 'ADMIN';
    insert into user_roles(user_id, role_id) values(v_user_id, v_role_id);

    select id into v_user_id from users where username = 'basilex';
    select id into v_role_id from roles where name = 'ADMIN';
    insert into user_roles(user_id, role_id) values(v_user_id, v_role_id);
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
    insert into profile(user_id, firstname, lastname) values(v_user_id, 'Superuser', 'Administrator');

    select id into v_user_id from users where username = 'admin';
    insert into profile(user_id, firstname, lastname) values(v_user_id, 'System', 'Administrator');

    select id into v_user_id from users where username = 'basilex';
    insert into profile(user_id, firstname, lastname, gender, birthday) values(v_user_id, 'Alexander', 'Vasilenko', 'male', '1965-04-03'::date);
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
