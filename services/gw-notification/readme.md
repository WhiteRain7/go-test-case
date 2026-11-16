# Project/gw-notification/

Gets huge transfer messages from wallets via kafka and store them in mongoDB.

# How to run

## Docker run

No additional steps required.

## Local run

No additional steps required.

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
          - `mongoserver/` - connection to mongoDB to store messages
  - `transfers/` - domain, repo, usecase, and handlers to store messages in mongoDB
- `logs/` - directory for logs, local only

# environment

Uses `broker` and `logger` from `pkg/`

## mongo

```dotenv
host=notification_mongo # docker container name or localhost
port=27017 # mongo port
user=admin # credentials
password=admin # credentials
base=main # mongo database to use
```

## analytics

```dotenv
port=8081 # http port provides /metrics request for prometheus
```