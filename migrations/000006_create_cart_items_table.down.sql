ALTER TABLE IF EXISTS cart_items
    DROP CONSTRAINT IF EXISTS fk_cart_sessions;

ALTER TABLE IF EXISTS cart_items
    DROP CONSTRAINT IF EXISTS fk_cart_products;

DROP TABLE IF EXISTS cart_items;