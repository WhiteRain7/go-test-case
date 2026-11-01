# GOLANG

FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./out/api ./cmd/api

# RUNTIME

FROM alpine:3.19
RUN adduser -D -u 10001 app
RUN mkdir -p /logs && chown app:app /logs
USER app
COPY --from=build /src/out/api /api
EXPOSE 8080

ENTRYPOINT [ "/api" ]