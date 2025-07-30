# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init -g cmd/server/main.go --parseInternal --output docs

RUN go build -o online-subscriptions cmd/server/main.go

# Stage 2: Final runtime image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/online-subscriptions .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/internal/migrations ./migrations
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./online-subscriptions"]
