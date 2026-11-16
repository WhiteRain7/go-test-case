# GOLANG

FROM golang:1.27-alpine AS build
WORKDIR /src

COPY --from=pkg . ./pkg/
COPY --from=proto . ./proto/
COPY . ./services/service/

RUN --mount=type=cache,target=/go/pkg/mod cd ./services/service && go mod download
RUN --mount=type=cache,target=/go/pkg/mod cd ./services/service && go test ./...
RUN --mount=type=cache,target=/go/pkg/mod cd ./services/service && CGO_ENABLED=0 GOOS=linux go build -o /src/out ./cmd/api

# RUNTIME

FROM alpine:3.19
RUN adduser -D -u 10001 app
RUN mkdir -p /logs && chown app:app /logs
USER app
COPY --from=build /src/out /exchanger
EXPOSE 8081

ENTRYPOINT [ "/exchanger" ]