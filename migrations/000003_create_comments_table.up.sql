CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    product_id bigint not null,
    owner_id bigint not null,
    text text not null,
    creation_date timestamp(0) with time zone not null default NOW(),
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE IF EXISTS comments
    ADD CONSTRAINT fk_comments_products FOREIGN KEY (product_id) REFERENCES products (id);

ALTER TABLE IF EXISTS comments
    ADD CONSTRAINT fk_comments_users FOREIGN KEY (owner_id) REFERENCES users (id);