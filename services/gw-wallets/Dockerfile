# GOLANG

FROM golang:1.27-alpine AS build
WORKDIR /src

COPY --from=pkg . ./pkg/
COPY --from=proto . ./proto/
COPY . ./services/service/

RUN --mount=type=cache,target=/go/pkg/mod cd ./services/service && go mod download
RUN --mount=type=cache,target=/go/pkg/mod cd ./services/service && go test ./...
RUN --mount=type=cache,target=/go/pkg/mod cd ./services/service && CGO_ENABLED=0 GOOS=linux go build -o /src/out ./cmd/api

RUN --mount=type=cache,target=/go/pkg/mod go install github.com/swaggo/swag/cmd/swag@latest
RUN cd ./services/service && mkdir docs && swag init -g cmd/api/main.go -o docs

# RUNTIME

FROM alpine:3.19
RUN adduser -D -u 10001 app
RUN mkdir -p /logs && chown app:app /logs
USER app
COPY --from=build /src/out /wallets
COPY --from=build /src/services/service/docs /docs
EXPOSE 8080

ENTRYPOINT [ "/wallets" ]