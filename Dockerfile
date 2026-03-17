# 1. Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
# -o main: Output binary name
# ./cmd/yenup/main.go: Entry point path
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/yenup/main.go

# 2. Runtime stage
FROM debian:bullseye-slim

# Install CA certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Cloud Run sets the PORT env variable automatically (default 8080)
EXPOSE 8080

# Run the binary
CMD ["./main"]
