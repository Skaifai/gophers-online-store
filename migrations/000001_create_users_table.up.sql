CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    role varchar(10) not null,
    name varchar(20) not null,
    email citext UNIQUE not null,
    password_hash bytea not null,
    registration_date timestamp(0) with time zone not null default NOW(),
    profile bigint UNIQUE not null,
    version integer NOT NULL DEFAULT 1
);