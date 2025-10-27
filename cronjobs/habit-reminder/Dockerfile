# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /habit-reminder-cronjob ./cronjobs/habit-reminder

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /habit-reminder-cronjob .
USER nobody:nobody
ENTRYPOINT ["/habit-reminder-cronjob"]
