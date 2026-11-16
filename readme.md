# Project

Used:
- golang 1.27
- postgres 18
- goose (db migrations)
- redis
- mongodb
- kafka
- clickhouse
- prometheus
- grafana

Requires Docker 2.20.0+ ("include" directive used)

# How to run

## Docker run

1. For each service:
   1. Copy .env.example to .env.prod (for docker)
   2. Fill in .env.prod with actual values
2. Run `docker compose up -d`
3. By default, main service runs on 8080

## Local run

1. For each service:
   1. Copy .env.example to .env.prod (for docker)
   2. Fill in .env.prod with actual values
   3. `go mod download`
   4. service readme steps
   5. `go run cmd/api/main.go`
2. By default, main service runs on 8080

# Default addresses

- `http://localhost:8080` - main API
- `http://localhost:8080/docs` - Swagger API Docs
- `http://localhost:3000` - grafana (runtime analytics), default env credentials `admin`:`admin`

# Project structure

- `pkg/` - golang modules and utils shared between services, imported as is
- `proto/` - gRPC proto files, including compiled
- `services/`
  - `gw-wallets/` - main public service, handles HTTP requests from users
  - `gw-exchanger/` - provides gw-wallets with actual rates by gRPC
  - `gw-notification/` - receives transfers from gw-wallets by kafka and store them at mongoDB
  - `gw-analytics/` - retrieves info from services with prometheus and displays it with grafana

See readme.md files in `pkg/` and `servies/[service]/` for more information.