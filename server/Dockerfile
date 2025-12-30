FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build API server
RUN CGO_ENABLED=0 go build -o api ./cmd/api

# Build Seed command
RUN CGO_ENABLED=0 go build -o seed ./cmd/seed

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Copy binaries
COPY --from=builder /app/api .
COPY --from=builder /app/seed .

# Copy migrations files (for auto-migrate on startup)
COPY --from=builder /app/internal/migrations ./migrations

# Copy challenges files (for seed)
COPY --from=builder /app/challenges ./challenges

EXPOSE 8080
CMD ["./api"]
