# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o my-redis ./cmd/my-redis

# Final stage
FROM scratch

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/my-redis .

# Expose the default Redis port
EXPOSE 6378

# Create volume for data persistence
VOLUME ["/data"]

# Run the binary
CMD ["./my-redis", "-addr", "0.0.0.0:6378"]