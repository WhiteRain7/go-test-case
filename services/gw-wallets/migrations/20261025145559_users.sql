-- +goose Up
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(255) UNIQUE,
    "password" VARCHAR(60) NOT NULL,
    "email" VARCHAR(255) UNIQUE
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_username" ON users(username);

-- +goose Down
DROP TABLE IF EXISTS "users";

DROP INDEX IF EXISTS "users_username";