-- +goose Up
CREATE TABLE IF NOT EXISTS "rates" (
    "currency" VARCHAR(3) PRIMARY KEY,
    "rate" FLOAT NOT NULL DEFAULT 0.0
);

INSERT INTO "rates" VALUES
('RUB', 1),
('EUR', 100),
('USD', 90);

-- +goose Down
DROP TABLE IF EXISTS "rates";
