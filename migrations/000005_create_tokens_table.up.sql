CREATE TABLE IF NOT EXISTS tokens (
    id bigserial PRIMARY KEY,
    user_id bigint not null,
    refresh_token text not null
);

ALTER TABLE IF EXISTS tokens
    ADD CONSTRAINT fk_tokens_users FOREIGN KEY (user_id) REFERENCES users (id);