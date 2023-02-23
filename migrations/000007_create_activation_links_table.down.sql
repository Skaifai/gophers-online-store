ALTER TABLE IF EXISTS activation_links
    DROP CONSTRAINT IF EXISTS fk_activations_users;

DROP TABLE IF EXISTS activation_links;