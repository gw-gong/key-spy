# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy vendor directory and go files
COPY vendor/ vendor/
COPY go.mod go.sum ./

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/
COPY config/ config/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/bin/scanner ./cmd/scanner

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config
COPY --from=builder /app/bin/scanner /app/scanner
COPY --from=builder /app/config /app/config

# Create output directory
RUN mkdir -p /data/key-spy/output /data/key-spy/logs

# Set timezone
ENV TZ=Asia/Shanghai

# Run the binary
ENTRYPOINT ["/app/scanner"]
CMD ["-cfg_path", "/app/config/scanner/localcfg", "-cfg_name", "live.yaml"]
