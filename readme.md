# Project

Used golang 1.27 and psql 18.

# How to run

## Docker run

1. Copy .env.example to .env.prod (for docker)
2. Fill in .env.prod with actual values
   - You may want to set `db.host` to `db` and `db.name` to `wallets` according to compose
3. Run `docker compose up -d`
4. By default, runs on 8080

## Local run

1. Copy .env.example to .env (local run)
2. Fill in .env with actual values:
    - For local run you may want to set goose env variables (example below)
3. `go mod download`
4. `goose up`
5. `go run cmd/api/main.go`
6. By default, runs on 8080

```dotenv
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="user=postgres password=postgres dbname=wallets sslmode=disable host=localhost port=5432"
GOOSE_MIGRATION_DIR=migrations
```

# Requests

## POST /wallets/create

`POST http://localhost:8080/api/v1/wallets/create`
```
body not needed
```

Returns new wallet with ID for other requests and 0 amount.

## GET /wallets/{uuid}

`GET http://localhost:8080/api/v1/wallets/{uuid}`

Returns wallet balance on success.

Returns 404 if uuid is not found.

## POST /wallet

`POST http://localhost:8080/api/v1/wallet`
```json
{
  "walletId": "{uuid}",
  "operation": "DEPOSIT or WITHDRAW",
  "amount": 100
}
```

Deposit adds amount to wallet balance. Withdraw subtracts amount from wallet balance.

Returns wallet with new balance on success.

Returns 404 if uuid is not found.

Returns 400, if operation is unknown (case-sensitive), amount is negative or
there's no enough funds to perform withdraw.

``

# Project structure

- `cmd/`
  - `api/`
    - `main.go` - entry point to application, just runs internal/app/app.go
- `internal/`
  - `app/`
    - `app.go` - initializes all resources and runs http service
  - `shared/`
    - `infrastructure/`
      - `apperr/apperr.go` - wrap for errors to provide HTTP error codes
      - `config/config.go` - parses .env file for configuration
      - `http/http.go` - connects handlers and runs service
      - `logger/logger.go` - async logger to file and console
      - `postgres/postgres.go` - connection to psql
  - `wallets/`
    - `domain/domain.go` - wallet struct with some methods
    - `repository/repository.go` - database handlers
    - `usecase/usecase.go` - use cases for business logic
    - `handler/handler.go` - http handlers
- `logs/` - directory for logs
- `migrations/` - database migrations, powered by goose