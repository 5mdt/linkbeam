# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-s -w -X main.Version=${VERSION:-dev} -X main.Commit=${COMMIT:-unknown}" \
    -o linkbeam \
    ./cmd/linkbeam

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/linkbeam /app/linkbeam

# Copy templates and themes
COPY templates /app/templates
COPY themes /app/themes

# Create config directory
RUN mkdir -p /app/config

# Expose port (if needed for future HTTP server)
EXPOSE 8080

# Set default config path
ENV CONFIG_PATH=/app/config/config.yaml

# Run as non-root user
RUN addgroup -g 1000 linkbeam && \
    adduser -D -u 1000 -G linkbeam linkbeam && \
    chown -R linkbeam:linkbeam /app

USER linkbeam

ENTRYPOINT ["/app/linkbeam"]
CMD ["--config", "/app/config/config.yaml", "--output", "/app/dist/index.html"]
