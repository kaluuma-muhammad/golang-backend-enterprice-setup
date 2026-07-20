# -----------------------------
# Stage 1 - Build
# -----------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o api \
    ./cmd/api

# Build Seeder
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o seed \
    ./cmd/seed

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# -----------------------------
# Stage 2 - Runtime
# -----------------------------
FROM alpine:3.22

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    postgresql17-client

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Binaries
COPY --from=builder /app/api .
COPY --from=builder /app/seed .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Database migrations
COPY internal/infrastructure/postgres/migrations ./migrations

# Startup script
COPY deployments/docker/entrypoint.sh .

RUN chmod +x entrypoint.sh

USER app

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]