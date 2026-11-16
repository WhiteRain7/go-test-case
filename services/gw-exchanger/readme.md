# Project/gw-exchanger/

Service provide actual rates to other services, more specifically - gw-wallets.

# How to run

## Docker run

No additional steps required.

## Local run

Insert goose env to .env file.

```dotenv
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="user=postgres password=postgres dbname=exchanger sslmode=disable host=localhost port=5432"
GOOSE_MIGRATION_DIR=migrations
```

Adapt for own use, create psql database, then run `goose up`.

# Project structure

- `cmd/`
  - `api/`
      - `main.go` - entry point to application, just runs internal/app/app.go
- `internal/`
  - `app/`
      - `app.go` - initializes all resources and runs grpc service
  - `shared/`
      - `infrastructure/`
          - `config/` - see `/pkg/infrastructure/config/`
          - `promanalytics/` - provides prometheus metrics
          - `protoserver/` - provides gRPC server
  - `rates/` - domain, repo, usecase, and handlers to get rates from db
  - `adapters/`
    - `grpc/` - grpc adapters for proto files
- `logs/` - directory for logs, local only
- `migrations/` - database migrations, docker volume, powered by goose

# environment

Uses `postgres` (with prefix db) and `logger` from `pkg/`

## protoserver

```dotenv
port=9001 # exposed port to listen
```

## analytics

```dotenv
port=8081 # http port provides /metrics request for prometheus
```