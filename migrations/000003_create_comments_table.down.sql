ALTER TABLE IF EXISTS comments
    DROP CONSTRAINT IF EXISTS fk_comments_products;

ALTER TABLE IF EXISTS comments
    DROP CONSTRAINT IF EXISTS fk_comments_users;

DROP TABLE IF EXISTS comments;