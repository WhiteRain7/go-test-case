-- +goose Up
CREATE TABLE IF NOT EXISTS analytics.transfers (
   user      UInt64,
   operation String,
   amount    Float64,
   currency  String,
   timestamp DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (timestamp, operation, amount);

-- +goose Down
DROP TABLE IF EXISTS analytics.transfers;
