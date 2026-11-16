-- +goose Up
CREATE TABLE IF NOT EXISTS "wallets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "amount" INT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS "wallets";
