CREATE TABLE IF NOT EXISTS cart_items (
    id bigserial PRIMARY KEY,
    session_id bigint not null,
    product_id bigint not null,
    quantity int not null,
    creation_date timestamp(0) with time zone not null default NOW()
);

ALTER TABLE IF EXISTS cart_items
    ADD CONSTRAINT fk_cart_sessions FOREIGN KEY (session_id) REFERENCES shopping_sessions (id);

ALTER TABLE IF EXISTS cart_items
    ADD CONSTRAINT fk_cart_products FOREIGN KEY (product_id) REFERENCES products (id);