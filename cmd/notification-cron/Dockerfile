# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-cron ./cmd/notification-cron

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests to Twilio
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/notification-cron .

# Use non-root user
USER appuser

# Run the application
ENTRYPOINT ["/app/notification-cron"]
