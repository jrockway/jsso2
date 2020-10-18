-- Write your migrate up statements here
create table "user" (id bigserial primary key not null, username text not null);
create unique index user_username on "user" (username);
