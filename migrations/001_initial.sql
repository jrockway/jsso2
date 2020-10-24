-- Write your migrate up statements here
create table "user" (id bigserial primary key not null, username text not null);
create unique index user_username on "user" (username);

create table session (
    id bytea primary key not null check (octet_length(id) = 64),
    user_id bigint not null,
    metadata jsonb not null,
    created_at timestamp (3) with time zone not null,
    expires_at timestamp (3) with time zone null,
    constraint fk_user foreign key (user_id) references "user" (id)
);
