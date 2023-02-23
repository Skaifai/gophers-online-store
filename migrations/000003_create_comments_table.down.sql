ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS fk_comments_products;

ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS fk_comments_users;

DROP TABLE IF EXISTS comments;