# Project/gw-analytics/

Retrieves various metrics from other services and displays it with grafana.

# How to run

## Docker run

Grafana should be configured manually or imported from json in this directory.

## Local run

Grafana should be configured manually or imported from json in this directory.

# Project structure

- `cmd/`
  - `api/`
      - `main.go` - entry point to application, just runs internal/app/app.go
- `internal/`
  - `app/`
      - `app.go` - initializes all resources and runs broker consumer service
  - `shared/`
      - `infrastructure/`
          - `config/` - see `/pkg/infrastructure/config/`
          - `promanalytics/` - provides prometheus metrics
          - `clickhouseserver/` - connection to clickhouse to store messages
  - `transfers/` - domain, repo, usecase, and handlers to store messages in clickhouse
- `logs/` - directory for logs, local only
- `migrations/` - database migrations, docker volume, powered by goose

# environment

Uses `broker` and `logger` from `pkg/`

## clickhouse

```dotenv
host=clickhouse # docker container name or localhost
port=9000 # clickhouse port
user=admin # credentials
password=admin # credentials
base=analytics # mongo database to use
```

## analytics

```dotenv
port=8081 # http port provides /metrics request for prometheus
```