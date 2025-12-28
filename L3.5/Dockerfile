# Build stage
FROM golang:1.25-alpine AS builder

# Set GOTOOLCHAIN to auto to allow downloading newer Go versions
ENV GOTOOLCHAIN=auto

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o event-booker ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/event-booker .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./event-booker"]

