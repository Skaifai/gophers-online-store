CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    role varchar(10) not null default 'USER',
    username varchar(20) UNIQUE not null,
    email citext UNIQUE not null,
    phone_number varchar(15) not null,
    password_hash bytea not null,
    registration_date timestamp(0) with time zone not null default NOW(),
    name varchar(20) not null,
    surname varchar(20) not null,
    date_of_birth date not null,
    address text not null,
    about_me text not null default 'Tell us about yourself.',
    picture_url text not null default 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png',
    activated boolean not null default false,
    version integer NOT NULL DEFAULT 1
);