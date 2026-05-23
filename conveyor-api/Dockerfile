FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY conveyor-api/go.mod conveyor-api/go.sum ./
RUN go mod download
COPY conveyor-api/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

RUN apk add --no-cache curl && \
    curl -sL https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz -o migrate.tar.gz && \
    tar xzf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    rm migrate.tar.gz

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /api /api
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/dashboard /dashboard
COPY --from=builder /app/migrations /migrations

COPY conveyor-api/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 8080
ENTRYPOINT ["/start.sh"]
