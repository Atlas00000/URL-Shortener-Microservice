# Build stage
FROM golang:1.22-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create necessary directories
RUN mkdir -p data/geoip

# Download GeoIP database
ARG MAXMIND_LICENSE_KEY
RUN curl -L "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz" -o GeoLite2-City.tar.gz && \
    tar -xzf GeoLite2-City.tar.gz && \
    mv GeoLite2-City_*/GeoLite2-City.mmdb data/geoip/ && \
    rm -rf GeoLite2-City_* GeoLite2-City.tar.gz

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o urlshortener ./src/main.go

# Final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the binary and data from builder
COPY --from=builder /app/urlshortener .
COPY --from=builder /app/data ./data

# Create necessary directories
RUN mkdir -p data/geoip

# Set environment variables
ENV DATA_DIR=/app/data
ENV PORT=8080

# Expose port
EXPOSE 8080

# Run the application
CMD ["./urlshortener"] 