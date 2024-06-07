# Build stage
FROM golang:1.22 as builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Copy the .env file
COPY .env ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . ./

# Update package lists and install ca-certificates
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM scratch

# Set the working directory in the final image
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy CA certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set the entrypoint to the built binary
ENTRYPOINT ["/app/main"]
