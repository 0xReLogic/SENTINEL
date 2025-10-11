# Use Docker BuildKit for better caching and multi-platform builds
# syntax=docker/dockerfile:1

FROM golang:1.21.13-alpine3.20 AS builder

# Build stage metadata
LABEL stage="builder" \
      maintainer="SENTINEL Team <team@sentinel.dev>" \
      version="1.0.0" \
      description="SENTINEL Builder Stage"

# Install build dependencies and create non-root user
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Set working directory for build
WORKDIR /build

# Copy Go module files first for better layer caching
COPY go.mod go.sum ./

# Download and verify Go dependencies (cached unless go.mod/go.sum changes)
RUN go mod download && \
    go mod verify

# Copy application source code
COPY . .

# Build static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s -X main.version=1.0.0 -extldflags '-static'" \
    -trimpath \
    -tags netgo \
    -o sentinel .

# Test the binary works before copying to runtime stage
RUN ./sentinel --help && ./sentinel validate

FROM alpine:3.20.3 AS runtime

# Build arguments for reproducible builds
ARG BUILD_DATE
ARG VERSION=1.0.0

# Runtime stage metadata following OCI standards
LABEL org.opencontainers.image.title="SENTINEL" \
      org.opencontainers.image.description="A simple monitoring system written in Go" \
      org.opencontainers.image.vendor="SENTINEL Team" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.source="https://github.com/0xReLogic/SENTINEL" \
      org.opencontainers.image.licenses="MIT"

# Install runtime dependencies and create secure user
RUN apk --no-cache add \
    ca-certificates \
    tzdata && \
    addgroup -g 1001 -S sentinel && \
    adduser -S sentinel -u 1001 -G sentinel -h /app && \
    mkdir -p /app/data && \
    chown -R sentinel:sentinel /app

# Set application working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/sentinel .

# Copy configuration file
COPY sentinel.yaml .

# Set proper ownership and read-only permissions for security
RUN chown sentinel:sentinel /app/sentinel /app/sentinel.yaml && \
    chmod 555 /app/sentinel && \
    chmod 444 /app/sentinel.yaml

# Run as non-root user for security
USER sentinel

# Health check using sentinel once command
HEALTHCHECK --interval=60s --timeout=15s --retries=3 --start-period=30s \
    CMD ["sh", "-c", "./sentinel once > /dev/null 2>&1 && echo 'healthy' || exit 1"]

# Expose port for future web dashboard feature
EXPOSE 8080

# Start sentinel in continuous monitoring mode
CMD ["./sentinel", "run"]
