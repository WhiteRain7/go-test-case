-- +goose Up
CREATE TABLE "currency_types" (
    "key" VARCHAR(3) PRIMARY KEY
);

INSERT INTO "currency_types" VALUES
('RUB'),
('EUR'),
('USD');

CREATE TABLE "wallet_currencies" (
    "wallet" UUID REFERENCES wallets(id),
    "currency" VARCHAR(3) NOT NULL REFERENCES currency_types(key),
    "amount" FLOAT DEFAULT 0,

    PRIMARY KEY ("wallet", "currency")
);

INSERT INTO "wallet_currencies"
SELECT "id", 'RUB', "amount"
FROM "wallets";

ALTER TABLE "wallets" DROP COLUMN "amount";
ALTER TABLE "wallets" ADD COLUMN "user" INTEGER REFERENCES users(id);

INSERT INTO "wallets" ("user")
SELECT "id"
FROM "users";

-- +goose Down
ALTER TABLE "wallets" ADD COLUMN "amount" INT NOT NULL DEFAULT 0;
ALTER TABLE "wallets" DROP COLUMN "user";

UPDATE "wallets"
SET "amount" = FLOOR(s."amount")::integer
FROM (
    SELECT "wallet", "amount"
    FROM "wallet_currencies"
    WHERE "currency" == 'RUB'
) AS s
WHERE "id" = s.wallet;

DROP TABLE "wallet_currencies";
DROP TABLE "currency_types";
