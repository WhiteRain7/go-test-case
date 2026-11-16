-- +goose Up
ALTER TABLE wallet_currencies
ALTER COLUMN amount TYPE INTEGER
    USING ROUND(amount * 100)::INTEGER;

-- +goose Down
ALTER TABLE wallet_currencies
ALTER COLUMN amount TYPE FLOAT
    USING (amount::FLOAT)/100;
