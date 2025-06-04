# Build stage
FROM golang:1.21-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o url-shortener ./src/main.go

# Final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Create data directory for GeoIP database
RUN mkdir -p /app/data/geoip

# Copy the binary from builder
COPY --from=builder /app/url-shortener .

# Copy static files and set permissions
COPY --from=builder /app/static /app/static
RUN chmod -R 755 /app/static

# Set environment variables
ENV DATA_DIR=/app/data
ENV PORT=8080
ENV REDIS_URL="redis://redis:6379/0"

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./url-shortener"] 