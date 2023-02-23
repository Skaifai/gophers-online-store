CREATE TABLE IF NOT EXISTS activation_links (
    link text PRIMARY KEY,
    activated boolean not null default false,
    user_id bigint not null,
    creation_date timestamp(0) not null default NOW(0)
);

ALTER TABLE IF EXISTS activation_links
    ADD CONSTRAINT fk_activations_users FOREIGN KEY (user_id) REFERENCES users (id);