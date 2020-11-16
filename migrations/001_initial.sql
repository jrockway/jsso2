-- Write your migrate up statements here
create extension citext;

create table "user" (
    id bigserial primary key not null,
    username citext not null
);
create unique index idx_unique_user_username on "user" (username);

create table session (
    id bytea primary key not null check (octet_length(id) = 64),
    user_id bigint not null,
    metadata jsonb not null check (metadata != 'null'),
    taints jsonb not null check (taints != 'null'),
    created_at timestamp (3) with time zone not null,
    expires_at timestamp (3) with time zone not null,
    constraint fk_user foreign key (user_id) references "user" (id)
);
create index idx_session_user on session (user_id);

create table credential (
    id bigserial primary key not null,
    credential_id bytea not null check (octet_length(credential_id) >= 16),
    public_key bytea not null,
    user_id bigint not null,
    name text not null,
    created_at timestamp (3) with time zone not null,
    deleted_at timestamp (3) with time zone null,
    created_by_session_id bytea not null,
    constraint fk_user foreign key (user_id) references "user" (id),
    constraint fk_session foreign key (created_by_session_id) references session (id)
);
create unique index idx_unique_active_credential on credential (credential_id, user_id, deleted_at);
create index idx_credential_user on credential (user_id);
create index idx_credential_user_active on credential (user_id) where deleted_at is null;
