CREATE TABLE IF NOT EXISTS products (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
     id bigserial PRIMARY KEY,
     name varchar(20) not null,
     price float8 not null,
     description text not null,
     category varchar(20) not null,
     quantity int not null default 0,
     is_available boolean not null default true,
     creation_date timestamp(0) with time zone not null default NOW(),
     version integer NOT NULL DEFAULT 1
);