ALTER TABLE IF EXISTS comments
    DROP CONSTRAINT IF EXISTS fk_sessions_users;

DROP TABLE IF EXISTS shopping_sessions;