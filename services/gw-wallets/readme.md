# Project/gw-wallets/

The main service that handles user requests by HTTP.

# How to run

## Docker run

No additional steps required.

## Local run

Insert goose env to .env file.

```dotenv
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="user=postgres password=postgres dbname=wallets sslmode=disable host=localhost port=5432"
GOOSE_MIGRATION_DIR=migrations
```

Adapt for own use, create psql database, then run `goose up`.

# Project structure

- `cmd/`
  - `api/`
      - `main.go` - entry point to application, just runs internal/app/app.go
- `internal/`
  - `app/`
      - `app.go` - initializes all resources and runs http service
  - `shared/`
      - `infrastructure/`
          - `config/` - see `/pkg/infrastructure/config/`
          - `auth/` - JWT authorization
          - `httpserver/` - connects handlers and runs service
          - `rediswraps/` - wraps repositories with redis cache
          - `promanalytics/` - provides prometheus metrics
          - `swagger/` - provides API docs
  - `wallets/` - domain, repo, usecase, and handlers for wallets, most operations
  - `users/` - domain, repo, usecase, and handlers for users, login and register
  - `rates/` - domain, repo, usecase, and handlers to get rates by gRPC from gw-exchanger
- `logs/` - directory for logs, local only
- `migrations/` - database migrations, docker volume, powered by goose

# environment

Uses `postgres` (with prefix db), `protoclient`, `logger`, `redis` and `broker` from `pkg/`

## http

```dotenv
port=8080 # exposed port to listen
```

## auth

```dotenv
jwtsecret=123 # any string used as JWT signature, keep it secret
lifetime=15 # authorization token lifetime in minutes
issuer=whiterain7app # JWT extra info, do not affect anything
```