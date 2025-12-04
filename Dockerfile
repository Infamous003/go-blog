FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copying go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copying everything else to /app
COPY . .

# Building the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# ---------- RUN STAGE ----------
FROM alpine:latest

WORKDIR /app

# INstalling CA certificates for HTTP calls
RUN apk add --no-cache ca-certificates

# Copying the binray from the builder stage
COPY --from=builder /app/api .

COPY --from=builder /app/migrations ./migrations

EXPOSE 9090

CMD ./api -db-dsn="$GOBLOG_DSN"
