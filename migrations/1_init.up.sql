-- Устанавливаем расширение pgcrypto для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email        TEXT    NOT NULL UNIQUE,
    pass_hash    BYTEA   NOT NULL,
    is_admin     BOOLEAN DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);