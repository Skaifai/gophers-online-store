CREATE TABLE IF NOT EXISTS shopping_sessions (
    id bigserial PRIMARY KEY,
    user_id bigint not null,
    total float4 not null,
    creation_date timestamp(0) with time zone not null default NOW(),
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE IF EXISTS shopping_sessions
    ADD CONSTRAINT fk_sessions_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;